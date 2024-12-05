package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var opRegex = regexp.MustCompile(`(mul\((\d{1,3}),(\d{1,3})\))|(do|don't)\(\)`)

func main() {
	p1 := 0
	p2 := 0
	enabled := true
	for _, line := range getInput() {
		ops := opRegex.FindAllStringSubmatch(line, -1)
		for _, op := range ops {
			if op[1] != "" {
				l, err := strconv.Atoi(op[2])
				if err != nil {
					log.Fatal(err)
				}
				r, err := strconv.Atoi(op[3])
				if err != nil {
					log.Fatal(err)
				}
				p1 += l * r
				if enabled {
					p2 += l * r
				}
			} else if op[4] == "do" {
				enabled = true
			} else if op[4] == "don't" {
				enabled = false
			}
		}
	}
	fmt.Println("Part 1: ", p1)
	fmt.Println("Part 2: ", p2)
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
