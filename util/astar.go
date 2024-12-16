package util

import (
	"errors"
	"math"
	"slices"
)

type Graph[N comparable] interface {
	Start() N
	IsGoal(N) bool
	Neighbors(N) (neighbors []N, costs []int)
	Heuristic(N) int
}

var NoRoute = errors.New("no route")

func AStar[N comparable](graph Graph[N]) (int, []N, error) {
	start := graph.Start()
	open := map[N]struct{}{
		start: {},
	}
	gScore := make(map[N]int)
	gScore[start] = 0
	fScore := make(map[N]int)
	fScore[start] = graph.Heuristic(start)
	cameFrom := make(map[N]N)

	reconstruct := func(n N) []N {
		path := []N{n}
		for {
			p, ok := cameFrom[n]
			if !ok {
				slices.Reverse(path)
				return path
			}
			path = append(path, p)
			n = p
		}
	}

	bestNext := func() N {
		bestScore := math.MaxInt
		bestNode := start
		for n := range open {
			if score, ok := fScore[n]; ok && score < bestScore {
				bestScore = score
				bestNode = n
			}
		}
		if bestScore == math.MaxInt {
			panic("no best")
		}
		delete(open, bestNode)
		return bestNode
	}

	for len(open) > 0 {
		current := bestNext()
		if graph.IsGoal(current) {
			return gScore[current], reconstruct(current), nil
		}
		neighbors, costs := graph.Neighbors(current)
		if len(neighbors) != len(costs) {
			panic("unequal neighbors and costs")
		}
		for i, neighbor := range neighbors {
			edgeCost := costs[i]
			score := gScore[current] + edgeCost
			if oldScore, ok := gScore[neighbor]; !ok || score < oldScore {
				cameFrom[neighbor] = current
				gScore[neighbor] = score
				fScore[neighbor] = score + graph.Heuristic(neighbor)
				open[neighbor] = struct{}{}
			}
		}
	}
	return -1, nil, NoRoute
}
