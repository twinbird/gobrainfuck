package main

import (
	"fmt"
)

func makeJumpTable(pg *Program) (map[int]int, error) {
	pc := 0
	pgSize := pg.Size()
	jt := make(map[int]int)

	for pc < pgSize {
		inst, err := pg.Fetch(pc)
		if err != nil {
			return nil, err
		}
		if inst == '[' {
			nesting := 1
			seek := pc

			for {
				if nesting == 0 {
					break
				}
				seek++
				if seek >= pgSize {
					break
				}
				if inst, err := pg.Fetch(seek); err != nil {
					return nil, err
				} else if inst == ']' {
					nesting--
				} else if inst == '[' {
					nesting++
				}
			}
			if nesting == 0 {
				jt[pc] = seek
				jt[seek] = pc
			} else {
				return nil, fmt.Errorf("unmatched '[' at pc=%d", pc)
			}
		}
		pc++
	}

	return jt, nil
}

func (ip *Interpreter) OptimizedRun(pg *Program) error {
	jumpTable, err := makeJumpTable(pg)
	if err != nil {
		return err
	}
	for ip.programCounter < pg.Size() {
		inst, err := pg.Fetch(ip.programCounter)
		if err != nil {
			return err
		}

		switch inst {
		case '>':
			ip.dataPointer++
		case '<':
			ip.dataPointer--
		case '+':
			ip.memory[ip.dataPointer]++
		case '-':
			ip.memory[ip.dataPointer]--
		case '.':
			fmt.Fprintf(ip.output, "%c", ip.memory[ip.dataPointer])
		case ',':
			buf := make([]byte, 1)
			_, err := ip.input.Read(buf)
			if err != nil {
				return err
			}
			ip.memory[ip.dataPointer] = buf[0]
		case '[':
			if ip.memory[ip.dataPointer] != 0 {
				break
			}
			ip.programCounter = jumpTable[ip.programCounter]
		case ']':
			if ip.memory[ip.dataPointer] == 0 {
				break
			}
			ip.programCounter = jumpTable[ip.programCounter]
		default:
			return fmt.Errorf("instruction '%c' is bad char.PC=%d",
				ip.programCounter, inst)
		}
		ip.programCounter++
	}
	return nil
}
