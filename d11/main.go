package main

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/spikecurtis/aoc2024/util"
)

const example = "125 17"

type memoKey struct {
	engraving int
	blinks    int
}

var memo = map[memoKey]int64{}

func main() {
	lines := util.GetInputLines()
	stones := parseInput(lines[0])
	//stones = parseInput(example)

	p1 := int64(0)
	for _, s := range stones {
		p1 += numAfterBlinks(memoKey{
			engraving: s,
			blinks:    25,
		})
	}
	fmt.Println("Part 1:", p1)
	p2 := int64(0)
	for _, s := range stones {
		p2 += numAfterBlinks(memoKey{
			engraving: s,
			blinks:    75,
		})
	}
	fmt.Println("Part 2:", p2)
}

func numAfterBlinks(k memoKey) (result int64) {
	n, ok := memo[k]
	if ok {
		return n
	}
	defer func() { memo[k] = result }()
	if k.blinks == 0 {
		return 1
	}
	blinks := k.blinks - 1
	if k.engraving == 0 {
		return numAfterBlinks(memoKey{
			engraving: 1,
			blinks:    blinks,
		})
	}
	d := digits(k.engraving)
	if len(d)%2 == 0 {
		h := len(d) / 2
		l := memoKey{
			engraving: engraving(d[:h]),
			blinks:    blinks,
		}
		r := memoKey{
			engraving: engraving(d[h:]),
			blinks:    blinks,
		}
		return numAfterBlinks(l) + numAfterBlinks(r)
	}
	return numAfterBlinks(memoKey{
		engraving: k.engraving * 2024,
		blinks:    blinks,
	})
}

const base = 10

func digits(s int) []int {

	out := make([]int, 0)
	for s != 0 {
		d := s % base
		out = append(out, d)
		s -= d
		s /= base
	}
	slices.Reverse(out)
	return out
}

func engraving(dig []int) int {
	m := 1
	e := 0
	for i := len(dig) - 1; i >= 0; i-- {
		e += dig[i] * m
		m *= base
	}
	return e
}

func parseInput(line string) []int {
	s := strings.Split(line, " ")
	out := make([]int, len(s))
	var err error
	for i := range s {
		out[i], err = strconv.Atoi(s[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	return out
}
