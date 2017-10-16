package main

import (
	"fmt"
)

func (ip *Interpreter) SimpleRun(pg *Program) error {
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
			nesting := 1
			savedPC := ip.programCounter

			for {
				if nesting == 0 {
					break
				}
				ip.programCounter++
				if ip.programCounter >= pg.Size() {
					break
				}
				if i, err := pg.Fetch(ip.programCounter); err != nil {
					return err
				} else if i == ']' {
					nesting--
				} else if i == '[' {
					nesting++
				}
			}

			if nesting == 0 {
				break
			} else {
				return fmt.Errorf("unmatched '[' at pc=%d", savedPC)
			}
		case ']':
			if ip.memory[ip.dataPointer] == 0 {
				break
			}
			nesting := 1
			savedPC := ip.programCounter

			for (nesting != 0) && (ip.programCounter > 0) {
				ip.programCounter--
				if i, err := pg.Fetch(ip.programCounter); err != nil {
					return err
				} else if i == ']' {
					nesting++
				} else if i == '[' {
					nesting--
				}
			}

			if nesting == 0 {
				break
			} else {
				return fmt.Errorf("unmatched '[' at pc=%d", savedPC)
			}
		default:
			return fmt.Errorf("instruction '%c' is bad char.PC=%d",
				ip.programCounter, inst)
		}
		ip.programCounter++
	}
	return nil
}
