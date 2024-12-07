package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type equation struct {
	test    int
	numbers []int
}

var memo = make(map[string]bool)

func (e equation) soluble(part int) (result bool) {
	s, ok := memo[e.String(part)]
	if ok {
		return s
	}
	defer func() {
		memo[e.String(part)] = result
		if part == 1 && result {
			// part 1 soluble implies part 2 soluble
			memo[e.String(2)] = true
		}
	}()
	if len(e.numbers) == 1 {
		return e.numbers[0] == e.test
	}
	// *
	nums := make([]int, 0, len(e.numbers))
	nums = append(nums, e.numbers[0]*e.numbers[1])
	nums = append(nums, e.numbers[2:]...)
	em := equation{
		test:    e.test,
		numbers: nums,
	}
	if em.soluble(part) {
		return true
	}
	// +
	nums = make([]int, 0, len(e.numbers))
	nums = append(nums, e.numbers[0]+e.numbers[1])
	nums = append(nums, e.numbers[2:]...)
	ep := equation{
		test:    e.test,
		numbers: nums,
	}
	if ep.soluble(part) {
		return true
	}
	if part == 1 {
		return false
	}
	// concat
	nums = make([]int, 0, len(e.numbers))
	cs := fmt.Sprintf("%d%d", e.numbers[0], e.numbers[1])
	c, err := strconv.Atoi(cs)
	if err != nil {
		log.Fatal(err)
	}
	nums = append(nums, c)
	nums = append(nums, e.numbers[2:]...)
	ec := equation{
		test:    e.test,
		numbers: nums,
	}
	return ec.soluble(part)
}

func (e equation) String(part int) string {
	b := strings.Builder{}
	for _, n := range e.numbers {
		b.WriteString(strconv.Itoa(n))
		b.WriteRune(',')
	}
	return fmt.Sprintf("%d %d: %s", part, e.test, b.String())
}

func main() {
	eqs := parseInput(getInput())
	p1 := 0
	p2 := 0
	for _, eq := range eqs {
		if eq.soluble(1) {
			p1 += eq.test
		}
		if eq.soluble(2) {
			p2 += eq.test
		}
	}
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
}

func parseInput(lines []string) []equation {
	out := make([]equation, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, " ")
		testS := strings.Trim(parts[0], ":")
		test, err := strconv.Atoi(testS)
		if err != nil {
			log.Fatal(err)
		}
		nums := make([]int, len(parts)-1)
		for _, part := range parts[1:] {
			n, err := strconv.Atoi(part)
			if err != nil {
				log.Fatal(err)
			}
			nums = append(nums, n)
		}
		out[i] = equation{test, nums}
	}
	return out
}

func getInput() []string {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}
