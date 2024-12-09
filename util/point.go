package util

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

type Boundary2D[T Real] struct {
	MinX, MinY, MaxX, MaxY T
}

func (b Boundary2D[T]) In(point Point2D[T]) bool {
	return point.X >= b.MinX &&
		point.X <= b.MaxX &&
		point.Y >= b.MinY &&
		point.Y <= b.MaxY
}
