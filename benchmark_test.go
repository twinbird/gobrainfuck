package main

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var (
	factorProgram     *Program
	mandelbrotProgram *Program
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	factf, err := os.Open("sample_code/factor.bf")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	factorProgram, err = parseFromReader(factf)
	if err != nil {
		log.Fatal(err)
	}
	mandelf, err := os.Open("sample_code/mandelbrot.bf")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	mandelbrotProgram, err = parseFromReader(mandelf)
	if err != nil {
		log.Fatal(mandelbrotProgram)
	}
}

func BenchmarkFactorSimpleRun(b *testing.B) {
	in := bytes.NewBufferString("179424691\n")
	out := new(bytes.Buffer)

	ip := NewInterpreter(MemorySize, in, out)
	if err := ip.SimpleRun(factorProgram); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func BenchmarkMandelbrotSimpleRun(b *testing.B) {
	in := bytes.NewBufferString("")
	out := new(bytes.Buffer)

	ip := NewInterpreter(MemorySize, in, out)
	if err := ip.SimpleRun(mandelbrotProgram); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func BenchmarkFactorOptimizedRun(b *testing.B) {
	in := bytes.NewBufferString("179424691\n")
	out := new(bytes.Buffer)

	ip := NewInterpreter(MemorySize, in, out)
	if err := ip.OptimizedRun(factorProgram); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func BenchmarkMandelbrotOptimizedRun(b *testing.B) {
	in := bytes.NewBufferString("")
	out := new(bytes.Buffer)

	ip := NewInterpreter(MemorySize, in, out)
	if err := ip.OptimizedRun(mandelbrotProgram); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func BenchmarkFactorOptimized2Run(b *testing.B) {
	in := bytes.NewBufferString("179424691\n")
	out := new(bytes.Buffer)

	ip := NewInterpreter(MemorySize, in, out)
	if err := ip.Optimized2Run(factorProgram); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func BenchmarkMandelbrotOptimized2Run(b *testing.B) {
	in := bytes.NewBufferString("")
	out := new(bytes.Buffer)

	ip := NewInterpreter(MemorySize, in, out)
	if err := ip.Optimized2Run(mandelbrotProgram); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
