package main

import (
	"fmt"
)

const (
	INVALID_OP = iota
	INC_PTR
	DEC_PTR
	INC_DATA
	DEC_DATA
	READ_STDIN
	WRITE_STDOUT
	JUMP_IF_DATA_ZERO
	JUMP_IF_DATA_NOT_ZERO
)

type IROperation struct {
	opKind      int
	repeatCount int
}

type IRProgram struct {
	instructions []*IROperation
}

func (pg *IRProgram) Size() int {
	return len(pg.instructions)
}

func (pg *IRProgram) Fetch(pointer int) (*IROperation, error) {
	if pg.Size() <= pointer {
		return nil, fmt.Errorf("instruction pointer overflow")
	}
	return pg.instructions[pointer], nil
}

func makeJumpTable2(pg *IRProgram) (map[int]int, error) {
	pc := 0
	pgSize := pg.Size()
	jt := make(map[int]int)

	for pc < pgSize {
		inst, err := pg.Fetch(pc)
		if err != nil {
			return nil, err
		}
		if inst.opKind == JUMP_IF_DATA_ZERO {
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
				} else if inst.opKind == JUMP_IF_DATA_NOT_ZERO {
					nesting--
				} else if inst.opKind == JUMP_IF_DATA_ZERO {
					nesting++
				}
			}
			if nesting == 0 {
				jt[pc] = seek
				jt[seek] = pc
			} else {
				return nil, fmt.Errorf("unmatched '['")
			}
		}
		pc++
	}

	return jt, nil
}

func charToOperation(c byte) int {
	switch c {
	case '>':
		return INC_PTR
	case '<':
		return DEC_PTR
	case '+':
		return INC_DATA
	case '-':
		return DEC_DATA
	case '.':
		return WRITE_STDOUT
	case ',':
		return READ_STDIN
	case '[':
		return JUMP_IF_DATA_ZERO
	case ']':
		return JUMP_IF_DATA_NOT_ZERO
	default:
		return INVALID_OP
	}
}

func String(pg *IRProgram) string {
	ret := ""
	for i := 0; i < len(pg.instructions); i++ {
		switch pg.instructions[i].opKind {
		case INC_PTR:
			ret += ">"
		case DEC_PTR:
			ret += "<"
		case INC_DATA:
			ret += "+"
		case DEC_DATA:
			ret += "-"
		case WRITE_STDOUT:
			ret += "."
		case READ_STDIN:
			ret += ","
		case JUMP_IF_DATA_ZERO:
			ret += "["
		case JUMP_IF_DATA_NOT_ZERO:
			ret += "]"
		default:
			ret += "?"
		}
	}
	return ret
}

func translateIR(pg *Program) (*IRProgram, error) {
	irInstructions := make([]*IROperation, 0)
	pc := 0
	pgSize := pg.Size()

	for pc < pgSize {
		inst, err := pg.Fetch(pc)
		if err != nil {
			return nil, err
		}
		opKind := charToOperation(inst)
		op := &IROperation{
			opKind:      opKind,
			repeatCount: 1,
		}
		for (pc + 1) < pgSize {
			if opKind == JUMP_IF_DATA_ZERO || opKind == JUMP_IF_DATA_NOT_ZERO {
				break
			}
			if i, err := pg.Fetch(pc + 1); err != nil {
				return nil, err
			} else if inst == i {
				op.repeatCount += 1
				pc++
			} else {
				break
			}
		}
		irInstructions = append(irInstructions, op)
		pc++
	}

	irprg := &IRProgram{}
	irprg.instructions = irInstructions
	return irprg, nil
}

func (ip *Interpreter) Optimized2Run(pg *Program) error {
	ir, err := translateIR(pg)
	if err != nil {
		return err
	}
	jumpTable, err := makeJumpTable2(ir)
	if err != nil {
		return err
	}
	for ip.programCounter < ir.Size() {
		inst, err := ir.Fetch(ip.programCounter)
		if err != nil {
			return err
		}

		switch inst.opKind {
		case INC_PTR:
			ip.dataPointer += inst.repeatCount
		case DEC_PTR:
			ip.dataPointer -= inst.repeatCount
		case INC_DATA:
			ip.memory[ip.dataPointer] += byte(inst.repeatCount)
		case DEC_DATA:
			ip.memory[ip.dataPointer] -= byte(inst.repeatCount)
		case WRITE_STDOUT:
			for i := 0; i < inst.repeatCount; i++ {
				fmt.Fprintf(ip.output, "%c", ip.memory[ip.dataPointer])
			}
		case READ_STDIN:
			buf := make([]byte, 1)
			for i := 0; i < inst.repeatCount; i++ {
				_, err := ip.input.Read(buf)
				if err != nil {
					return err
				}
				ip.memory[ip.dataPointer] = buf[0]
			}
		case JUMP_IF_DATA_ZERO:
			if ip.memory[ip.dataPointer] != 0 {
				break
			}
			ip.programCounter = jumpTable[ip.programCounter]
		case JUMP_IF_DATA_NOT_ZERO:
			if ip.memory[ip.dataPointer] == 0 {
				break
			}
			ip.programCounter = jumpTable[ip.programCounter]
		default:
			return fmt.Errorf("instruction '%c' is bad char",
				ip.programCounter, inst)
		}
		ip.programCounter++
	}
	return nil
}
