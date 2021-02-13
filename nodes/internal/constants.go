package internal

import "math"

const (
	TAU = math.Pi * 2
)

func squareMag(p1, p2 *point) float64 {
	return math.Pow(p1.x - p2.x, 2) + math.Pow(p1.y - p2.y, 2)
}

func lerp(x, min, max float64) float64 {
	return (x * (max - min)) + min
}