package main

import (
	"syscall/js"

	"github.com/ojrac/opensimplex-go"
)

var noise *Noise

type Noise struct {
	cols      float64
	rows      float64
	threshold float64
	source    opensimplex.Noise
}

func NewNoise(seed int64, threshold, cols, rows float64) *Noise {
	return &Noise{
		source:    opensimplex.NewNormalized(seed),
		threshold: threshold,
		cols:      cols,
		rows:      rows,
	}
}

func (n *Noise) Eval(x, y, z float64) float64 {
	return n.source.Eval3(x, y, z)
}

func (n *Noise) getBit(points []interface{}, row, col int) int {
	if row > len(points)-1 {
		return 0
	} else if p := points[row].(js.Value).Index(col); p.IsUndefined() || p.Float() < n.threshold {
		return 0
	} else {
		return 1
	}
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("initNoise", js.FuncOf(initNoise))
	js.Global().Set("noise", js.FuncOf(r))
	<-c
}

func initNoise(this js.Value, args []js.Value) interface{} {
	var (
		seed      = args[0].Int()
		threshold = args[1].Float()
		cols      = args[2].Float()
		rows      = args[3].Float()
	)
	noise = NewNoise(int64(seed), threshold, cols, rows)
	return nil
}

func r(this js.Value, args []js.Value) interface{} {
	z := args[0].Float()
	mod := args[1].Float()
	points := make([]interface{}, 0)
	for row := 0.0; row < noise.rows; row++ {
		rPoints := make([]interface{}, 0)
		for col := 0.0; col < noise.cols; col++ {
			rPoints = append(rPoints, noise.Eval(col/mod, row/mod, z))
		}
		points = append(points, js.ValueOf(rPoints))
	}
	keys := make([]interface{}, 0)
	for row := 0; row < int(noise.rows); row++ {
		rKeys := make([]interface{}, 0)
		for col := 0; col < int(noise.cols); col++ {
			val := noise.getBit(points, row, col)<<3 | noise.getBit(points, row, col+1)<<2 | noise.getBit(points, row+1, col+1)<<1 | noise.getBit(points, row+1, col)
			rKeys = append(rKeys, val)
		}
		keys = append(keys, js.ValueOf(rKeys))
	}
	return js.ValueOf(map[string]interface{}{
		"points": points,
		"keys":   keys,
	})
}
