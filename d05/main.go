package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	predecessors, updates := parseInput(getInput())
	p1 := 0
	p2 := 0
	cmp := func(a, b int) int {
		if slices.Contains(predecessors[b], a) {
			return -1
		}
		return 1
	}
	for _, update := range updates {
		if slices.IsSortedFunc(update, cmp) {
			p1 += middle(update)
			continue
		}
		slices.SortFunc(update, cmp)
		p2 += middle(update)
	}
	fmt.Println("Part 1: ", p1)
	fmt.Println("Part 2: ", p2)
}

func middle(update []int) int {
	return update[len(update)/2]
}

func parseInput(lines []string) (predecessors map[int][]int, updates [][]int) {
	predecessors = make(map[int][]int)
	i := 0
	for {
		line := lines[i]
		i++
		if line == "" {
			break
		}
		parts := strings.Split(line, "|")
		pre, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		post, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal(err)
		}
		predecessors[post] = append(predecessors[post], pre)
	}
	for ; i < len(lines); i++ {
		line := lines[i]
		parts := strings.Split(line, ",")
		update := make([]int, len(parts))
		for j, part := range parts {
			page, err := strconv.Atoi(part)
			if err != nil {
				log.Fatal(err)
			}
			update[j] = page
		}
		updates = append(updates, update)
	}
	return predecessors, updates
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

func getExampleInput() []string {
	example := `47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`
	return strings.Split(example, "\n")
}
