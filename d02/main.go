package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var reports []report
	lines := getInput()
	for _, line := range lines {
		reports = append(reports, parseInput(line))
	}
	numSafe := 0
	for _, r := range reports {
		if r.safe() {
			numSafe++
		}
	}
	fmt.Println("Part 1: ", numSafe)
	numSafe = 0
	for _, r := range reports {
		if r.safe() {
			numSafe++
			continue
		}
		for j := range r {
			r2 := make(report, 0, len(r)-1)
			r2 = append(r2, r[:j]...)
			r2 = append(r2, r[j+1:]...)
			if r2.safe() {
				numSafe++
				break
			}
		}
	}
	fmt.Println("Part 2: ", numSafe)
}

type report []int

func parseInput(l string) report {
	parts := strings.Split(l, " ")
	levels := make(report, len(parts))
	for i, part := range parts {
		level, err := strconv.Atoi(part)
		if err != nil {
			log.Fatal(err)
		}
		levels[i] = level
	}
	return levels
}

func (r report) safe() bool {
	var sign int
	for i := range r {
		var d int
		if i == 0 {
			continue
		}
		d = r[i] - r[i-1]
		if d == 0 {
			return false
		}
		if i == 1 {
			if d < 0 {
				sign = -1
			} else {
				sign = 1
			}
		} else {
			if d < 0 && sign == 1 {
				return false
			}
			if d > 0 && sign == -1 {
				return false
			}
		}
		d *= sign
		if d < 1 || d > 3 {
			return false
		}
	}
	return true
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
