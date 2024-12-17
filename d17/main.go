package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type computer struct {
	A, B, C int
	pc      int
	program []int
}

var halt = errors.New("halt")

func (c *computer) run() []int {
	out := make([]int, 0)
	for {
		o, err := c.runOne()
		if err != nil {
			return out
		}
		out = append(out, o...)
	}
}

func (c *computer) runOne() ([]int, error) {
	if c.pc >= len(c.program) {
		return nil, halt
	}

	advance := true
	defer func() {
		if advance {
			c.pc += 2
		}
	}()

	var out []int
	opcode := c.program[c.pc]
	operand := c.program[c.pc+1]
	switch opcode {
	case 0: // adv
		c.A = c.A >> c.comboOperand(operand)
	case 1: // bxl
		c.B = c.B ^ operand
	case 2: // bst
		c.B = c.comboOperand(operand) % 8
	case 3: // jnz
		if c.A == 0 {
			break
		}
		advance = false
		c.pc = operand
	case 4: // bxc
		c.B = c.B ^ c.C
	case 5: // out
		out = append(out, c.comboOperand(operand)%8)
	case 6: // bdv
		c.B = c.A >> c.comboOperand(operand)
	case 7: // cdv
		c.C = c.A >> c.comboOperand(operand)
	}
	return out, nil
}

func (c *computer) comboOperand(operand int) int {
	switch {
	case operand < 4:
		return operand
	case operand == 4:
		return c.A
	case operand == 5:
		return c.B
	case operand == 6:
		return c.C
	default:
		panic(operand)
	}
}

func main() {
	c := realInput()
	out := c.run()
	outS := make([]string, len(out))
	for i, n := range out {
		outS[i] = fmt.Sprintf("%d", n)
	}
	fmt.Println("Part 1:", strings.Join(outS, ","))

	solvedA, ok := findA(c.program, 0, 0)
	if ok {
		fmt.Println("Part 2: ", solvedA)
		return
	}
	fmt.Println("no part 2 found")
}

func example() *computer {
	return &computer{
		A:       729,
		program: []int{0, 1, 5, 4, 3, 0},
	}
}

func realInput() *computer {
	return &computer{
		A:       47792830,
		program: []int{2, 4, 1, 5, 7, 5, 1, 6, 4, 3, 5, 5, 0, 3, 3, 0},
	}
}

func findA(program []int, A, bits int) (int, bool) {
	nextBits := bits + 3
	l := nextBits / 3
	if len(program) < l {
		return -1, false
	}
	for suffix := range 1 << 3 {
		test := A << 3
		test += suffix

		c := computer{
			A:       test,
			program: program,
		}

		out := c.run()
		if slices.Equal(program, out) {
			return test, true
		}
		if len(out) < l {
			continue
		}
		os := out[len(out)-l:]
		ps := program[len(program)-l:]
		if slices.Equal(os, ps) {
			fmt.Println(A, suffix, test, out, l)
			solvedA, ok := findA(program, test, nextBits)
			if ok {
				return solvedA, true
			}
		}
	}
	return -1, false
}
