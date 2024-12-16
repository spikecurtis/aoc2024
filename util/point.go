package util

import (
	"errors"
	"iter"
)

type Real interface {
	int | int64 | float64
}

type Point2D[T Real] struct {
	X, Y T
}

func (p Point2D[T]) Plus(other Point2D[T]) Point2D[T] {
	return Point2D[T]{p.X + other.X, p.Y + other.Y}
}

func (p Point2D[T]) Minus(other Point2D[T]) Point2D[T] {
	return Point2D[T]{p.X - other.X, p.Y - other.Y}
}

func (p Point2D[T]) Scale(a T) Point2D[T] {
	return Point2D[T]{p.X * a, p.Y * a}
}

func North[T Real]() Point2D[T] {
	return Point2D[T]{0, -1}
}

func South[T Real]() Point2D[T] {
	return Point2D[T]{0, 1}
}

func East[T Real]() Point2D[T] {
	return Point2D[T]{1, 0}
}

func West[T Real]() Point2D[T] {
	return Point2D[T]{-1, 0}
}

func NorthWest[T Real]() Point2D[T] {
	return North[T]().Plus(West[T]())
}

func NorthEast[T Real]() Point2D[T] {
	return North[T]().Plus(East[T]())
}

func SouthEast[T Real]() Point2D[T] {
	return South[T]().Plus(East[T]())
}

func SouthWest[T Real]() Point2D[T] {
	return South[T]().Plus(West[T]())
}

func (p Point2D[T]) CardinalNeighbors() []Point2D[T] {
	out := make([]Point2D[T], 4)
	out[0] = p.Plus(North[T]())
	out[1] = p.Plus(East[T]())
	out[2] = p.Plus(South[T]())
	out[3] = p.Plus(West[T]())
	return out
}

func Abs[T Real](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func (p Point2D[T]) L1Norm() T {
	return Abs(p.X) + Abs(p.Y)
}

type Boundary2D[T Real] struct {
	MinX, MinY, MaxX, MaxY T
}

func (b Boundary2D[T]) In(point Point2D[T]) bool {
	return point.X >= b.MinX &&
		point.X <= b.MaxX &&
		point.Y >= b.MinY &&
		point.Y <= b.MaxY
}

func (b Boundary2D[int]) Iterate() iter.Seq[Point2D[int]] {
	return func(yield func(Point2D[int]) bool) {
		for y := b.MinY; y <= b.MaxY; y++ {
			for x := b.MinX; x <= b.MaxX; x++ {
				if !yield(Point2D[int]{X: x, Y: y}) {
					return
				}
			}
		}
	}
}

type Grid2D[V any] struct {
	Values [][]V
	Bounds Boundary2D[int]
}

var OutOfBounds = errors.New("out of bounds")

func (g Grid2D[V]) At(p Point2D[int]) (v V, err error) {
	if !g.Bounds.In(p) {
		return v, OutOfBounds
	}
	return g.Values[p.Y][p.X], nil
}

func (g Grid2D[V]) Set(p Point2D[int], v V) error {
	if !g.Bounds.In(p) {
		return OutOfBounds
	}
	g.Values[p.Y][p.X] = v
	return nil
}

func (g Grid2D[V]) Iterate() iter.Seq2[Point2D[int], V] {
	return func(yield func(p Point2D[int], v V) bool) {
		for pnt := range g.Bounds.Iterate() {
			if !yield(pnt, g.Values[pnt.Y][pnt.X]) {
				return
			}
		}
	}
}
