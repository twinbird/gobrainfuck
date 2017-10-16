package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	MemorySize = 30000
)

type Program struct {
	instructions []byte
}

func (pg *Program) Size() int {
	return len(pg.instructions)
}

func (pg *Program) Fetch(pointer int) (byte, error) {
	if pg.Size() <= pointer {
		return 0, fmt.Errorf("instruction pointer overflow")
	}
	return pg.instructions[pointer], nil
}

func (pg *Program) String() string {
	return string(pg.instructions)
}

func parseFromReader(reader io.Reader) (*Program, error) {
	pg := &Program{instructions: make([]byte, 0)}
	buf := make([]byte, 1)

	for {
		_, err := reader.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch buf[0] {
		case '>':
			fallthrough
		case '<':
			fallthrough
		case '+':
			fallthrough
		case '-':
			fallthrough
		case '.':
			fallthrough
		case ',':
			fallthrough
		case '[':
			fallthrough
		case ']':
			pg.instructions = append(pg.instructions, buf[0])
		}
	}
	return pg, nil
}

type Interpreter struct {
	memory         []byte
	programCounter int
	dataPointer    int
	input          io.Reader
	output         io.Writer
}

func NewInterpreter(msize int, in io.Reader, out io.Writer) *Interpreter {
	ip := &Interpreter{}
	ip.memory = make([]byte, msize)
	ip.input = in
	ip.output = out
	return ip
}

func main() {
	var src io.Reader
	var srcfile string
	var optimize bool
	var optimize2 bool

	flag.StringVar(&srcfile, "file", "", "source file")
	flag.BoolVar(&optimize, "o", false, "Optimize option1")
	flag.BoolVar(&optimize2, "o2", false, "Optimize option2")
	flag.Parse()

	if len(srcfile) > 0 {
		f, err := os.Open(srcfile)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		src = f
	} else {
		src = os.Stdin
	}
	pg, err := parseFromReader(src)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	ip := NewInterpreter(MemorySize, os.Stdin, os.Stdout)

	if optimize == true {
		if err := ip.OptimizedRun(pg); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else if optimize2 == true {
		if err := ip.Optimized2Run(pg); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		if err := ip.SimpleRun(pg); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
