package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spikecurtis/aoc2024/util"
)

type machine struct {
	ax, ay, bx, by, px, py int
}

var noSolution = errors.New("no solution")

func (m machine) bestCost(part int) (int, error) {
	det := m.ax*m.by - m.ay*m.bx
	if det == 0 {
		// degenerate --- technically could still work in special case that
		// prize is in line with just one button, but this doesn't appear to
		// be the case for any inputs.
		return 0, noSolution
	}
	a := (m.by * m.px / det) - (m.bx * m.py / det)
	b := -(m.ay * m.px / det) + (m.ax * m.py / det)
	// check if solution is non-integral
	if a*m.ax+b*m.bx != m.px {
		return 0, noSolution
	}
	if a*m.ay+b*m.by != m.py {
		return 0, noSolution
	}
	// can't press negative times
	if a < 0 || b < 0 {
		return 0, noSolution
	}
	// part 1 restricts to 100 presses
	if part == 1 && (a > 100 || b > 100) {
		return 0, noSolution
	}
	return a*3 + b, nil
}

func main() {
	machines := parseInput(util.GetInputLines())
	//machines := parseInput(getExampleLines())
	p1 := 0
	p2 := 0
	for _, m := range machines {
		cost, err := m.bestCost(1)
		if err == nil {
			p1 += cost
		}
		m2 := m
		m2.px += 10000000000000
		m2.py += 10000000000000
		cost, err = m2.bestCost(2)
		if err == nil {
			p2 += cost
		}
	}
	fmt.Println(p1)
	fmt.Println(p2)
}

func parseInput(lines []string) []machine {
	var out []machine
	var this machine
	var err error
	for i, line := range lines {
		switch i % 4 {
		case 0: // Button A: X+78, Y+81
			parts := strings.Split(line, " ")
			this.ax, err = strconv.Atoi(strings.Trim(parts[2], "X+,"))
			if err != nil {
				log.Fatal(err)
			}
			this.ay, err = strconv.Atoi(strings.Trim(parts[3], "Y+"))
			if err != nil {
				log.Fatal(err)
			}
		case 1: // Button B: X+18, Y+92
			parts := strings.Split(line, " ")
			this.bx, err = strconv.Atoi(strings.Trim(parts[2], "X+,"))
			if err != nil {
				log.Fatal(err)
			}
			this.by, err = strconv.Atoi(strings.Trim(parts[3], "Y+"))
			if err != nil {
				log.Fatal(err)
			}
		case 2: //Prize: X=4806, Y=5504
			parts := strings.Split(line, " ")
			this.px, err = strconv.Atoi(strings.Trim(parts[1], "X=,"))
			if err != nil {
				log.Fatal(err)
			}
			this.py, err = strconv.Atoi(strings.Trim(parts[2], "Y="))
			if err != nil {
				log.Fatal(err)
			}
		case 3:
			out = append(out, this)
		}
	}
	return out
}

func getExampleLines() []string {
	return strings.Split(`Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279
`, "\n")
}
