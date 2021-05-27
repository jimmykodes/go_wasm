package internal

import "math"

const (
	TAU = math.Pi * 2
)

type Point interface {
	X() float64
	Y() float64
	Mag() float64
}

var Origin = origin{}

type origin struct{}

func (o origin) X() float64 {
	return 0
}
func (o origin) Y() float64 {
	return 0
}
func (o origin) Mag() float64 {
	return 0
}

// squareMag returns the square magnitude of the distance between 2 Point locations
// useful for figuring out relative distances between objects. Saves the computational
// expense of square roots
func squareMag(p1, p2 Point) float64 {
	return math.Pow(p1.X()-p2.X(), 2) + math.Pow(p1.Y()-p2.Y(), 2)
}

// distance returns the actual distance between 2 Point locations
func distance(p1, p2 Point) float64 {
	return math.Sqrt(squareMag(p1, p2))
}

// lerp will a point between 0 and 1 and redistribute it between min and max.
// ex: lerp(0.5, 0.5, 1) == 0.75 // since 0.75 is half way between .5 and 1
// ex: lerp(0.25, 12, 20) == 14 // 14 is .25 of the way between 12 and 20
func lerp(x, min, max float64) float64 {
	return (x * (max - min)) + min
}

// m3 takes in magnitudes of 2 vectors and the radian of their interior angle to return
// the magnitude of the new vector after adding the two vectors together
func m3(m1, m2, t1, t2 float64) float64 {
	alpha := math.Pi - t1 + t2
	return math.Sqrt(math.Pow(m1, 2) + math.Pow(m2, 2) - (2 * m1 * m2 * math.Cos(alpha)))
}
func t3(m1, m2, m3, t1, t2 float64) float64 {
	num := math.Pow(m1, 2) + math.Pow(m3, 2) - math.Pow(m2, 2)
	den := 2 * m1 * m3
	t3 := math.Acos(num / den)
	if math.Mod(TAU-t1+t2, TAU) > math.Pi {
		t3 = TAU - t3
	}
	return t3 + t1
}
