package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/spikecurtis/aoc2024/util"
)

type warehouse struct {
	grid  util.Grid2D[byte]
	robot util.Point2D[int]
}

func (w *warehouse) moveRobot(d util.Point2D[int]) {
	if c, err := w.grid.At(w.robot); err != nil || c != '@' {
		panic("badness")
	}
	if w.tryMove(w.robot, d, false) {
		w.robot = w.robot.Plus(d)
	}
}

func (w *warehouse) tryMove(origin, d util.Point2D[int], dryRun bool) bool {
	targets := []util.Point2D[int]{origin.Plus(d)}
	oc, err := w.grid.At(origin)
	if err != nil {
		log.Fatal(err)
	}
	// vertical pushes need to evaluate additional targets
	if d == util.North[int]() || d == util.South[int]() {
		switch oc {
		case '[':
			targets = append(targets, targets[0].Plus(util.East[int]()))
		case ']':
			targets = append(targets, targets[0].Plus(util.West[int]()))
		}
	}
	for _, target := range targets {
		tc, err := w.grid.At(target)
		if err != nil {
			log.Fatal(err)
		}
		if tc == '#' {
			return false
		}
		if tc == '.' {
			continue
		}
		// must be a barrel of some kind, see if we can move it.
		if !w.tryMove(target, d, true) {
			return false
		}
	}
	if dryRun {
		return true
	}
	// we can move!
	for _, target := range targets {
		tc, err := w.grid.At(target)
		if err != nil {
			log.Fatal(err)
		}
		if tc == '.' {
			continue
		}
		if !w.tryMove(target, d, false) {
			log.Fatal("failed to move")
		}
	}

	w.move(origin, d)
	if d == util.North[int]() || d == util.South[int]() {
		switch oc {
		case '[':
			w.move(origin.Plus(util.East[int]()), d)
		case ']':
			w.move(origin.Plus(util.West[int]()), d)
		}
	}
	return true
}

// move actually copies bytes to move stuff, and logs fatal error if stuff is
// in the way.
func (w *warehouse) move(origin, d util.Point2D[int]) {
	contents, err := w.grid.At(origin)
	if err != nil {
		log.Fatal(err)
	}
	target := origin.Plus(d)
	tc, err := w.grid.At(target)
	if err != nil {
		log.Fatal(err)
	}
	if tc != '.' {
		panic("not empty")
	}
	err = w.grid.Set(target, contents)
	if err != nil {
		log.Fatal(err)
	}
	err = w.grid.Set(origin, '.')
	if err != nil {
		log.Fatal(err)
	}
}

func (w *warehouse) print() {
	for _, line := range w.grid.Values {
		fmt.Println(string(line))
	}
}

func parseInput(lines []string) (w *warehouse, moves []util.Point2D[int]) {
	i := 0
	boxAndWalls := make([][]byte, 0)
	var robot util.Point2D[int]
	for ; len(lines[i]) != 0; i++ {
		line := []byte(lines[i])
		if rx := slices.Index(line, '@'); rx != -1 {
			robot.X = rx
			robot.Y = i
		}
		boxAndWalls = append(boxAndWalls, line)
	}
	moves = make([]util.Point2D[int], 0)
	for ; i < len(lines); i++ {
		for _, b := range []byte(lines[i]) {
			switch b {
			case '^':
				moves = append(moves, util.North[int]())
			case '>':
				moves = append(moves, util.East[int]())
			case '<':
				moves = append(moves, util.West[int]())
			case 'v':
				moves = append(moves, util.South[int]())
			}
		}
	}
	return &warehouse{
		grid: util.Grid2D[byte]{
			Values: boxAndWalls,
			Bounds: util.Boundary2D[int]{
				MinX: 0,
				MinY: 0,
				MaxX: len(boxAndWalls[0]) - 1,
				MaxY: len(boxAndWalls) - 1,
			},
		},
		robot: robot,
	}, moves
}

func (w *warehouse) widen() {
	boxAndWalls := make([][]byte, len(w.grid.Values))
	for i, old := range w.grid.Values {
		wide := make([]byte, 0, 2*len(old))
		for _, b := range old {
			switch b {
			case '#':
				wide = append(wide, '#', '#')
			case 'O':
				wide = append(wide, '[', ']')
			case '.':
				wide = append(wide, '.', '.')
			case '@':
				wide = append(wide, '@', '.')
			}
		}
		boxAndWalls[i] = wide
	}
	w.grid = util.Grid2D[byte]{
		Values: boxAndWalls,
		Bounds: util.Boundary2D[int]{
			MinX: 0,
			MinY: 0,
			MaxX: len(boxAndWalls[0]) - 1,
			MaxY: len(boxAndWalls) - 1,
		},
	}
	w.robot.X = w.robot.X * 2
}

func main() {
	//w, moves := parseInput(getExampleLines())
	w, moves := parseInput(util.GetInputLines())
	//w.print()
	//fmt.Println("------------------")
	for _, move := range moves {
		//fmt.Println(move)
		w.moveRobot(move)
		//w.print()
		//fmt.Println("------------------")
	}
	p1 := 0
	for pnt, contents := range w.grid.Iterate() {
		if contents == 'O' {
			gps := pnt.Y*100 + pnt.X
			p1 += gps
		}
	}
	fmt.Println("Part 1: ", p1)
	w, moves = parseInput(util.GetInputLines())
	w.widen()
	//w.print()
	//fmt.Println("------------------")
	for _, move := range moves {
		//fmt.Println(move)
		w.moveRobot(move)
		//w.print()
		//fmt.Println("------------------")
	}
	p2 := 0
	for pnt, contents := range w.grid.Iterate() {
		if contents == '[' {
			gps := pnt.Y*100 + pnt.X
			p2 += gps
		}
	}
	fmt.Println("Part 2: ", p2)
}

func getExampleLines() []string {
	return strings.Split(`########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

<^^>>>vv<v>>v<<`, "\n")
}

func getWideExampleLines() []string {
	return strings.Split(`#######
#...#.#
#.....#
#..OO@#
#..O..#
#.....#
#######

<vv<<^^<<^^`, "\n")
}
