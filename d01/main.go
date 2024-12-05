package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	lines := getInput()
	var left []int
	var right []int
	for _, line := range lines {
		p := strings.Split(line, " ")
		l, err := strconv.Atoi(p[0])
		if err != nil {
			log.Fatal(err)
		}
		left = append(left, l)
		r, err := strconv.Atoi(p[3])
		if err != nil {
			log.Fatal(err)
		}
		right = append(right, r)
	}
	sort.Ints(left)
	sort.Ints(right)
	d := 0
	for i := range len(left) {
		diff := left[i] - right[i]
		if diff < 0 {
			diff = -diff
		}
		d += diff
	}
	fmt.Println("Part 1: ", d)
	// Part 2
	counts := make(map[int]int)
	for _, r := range right {
		counts[r]++
	}
	s := 0
	for _, l := range left {
		s += counts[l] * l
	}
	fmt.Println("Part 2: ", s)
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
