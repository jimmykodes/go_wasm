package internal

// Well is a source or a sink, either pushing points away from, or pulling them to, their x/y coords
// a source will have a positive strength and a sink will have a negative strength
type Well struct {
	x         float64
	y         float64
	strength  float64
	squareMag float64
}

func NewWell(x float64, y float64, strength float64) *Well {
	well := &Well{x: x, y: y, strength: strength,}
	well.squareMag = squareMag(Origin, well)
	return well
}

func (w Well) X() float64 {
	return w.x
}

func (w Well) Y() float64 {
	return w.y
}

func (w Well) Mag() float64 {
	return w.strength
}
