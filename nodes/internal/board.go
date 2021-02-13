package internal

import (
	"math"
	"math/rand"
	"sync"
	"syscall/js"
	"time"
)

type Board struct {
	initCalled   bool
	width        float64
	height       float64
	kLines       int
	points       []*point
	bounceBounds bool
	threshold    float64
}

func (b *Board) Init(this js.Value, args []js.Value) interface{} {
	if b.initCalled {
		return nil
	}
	var (
		width        = args[0].Float()
		height       = args[1].Float()
		nPoints      = args[2].Int()
		bounceBounds = args[3].Bool()
		kLines       = args[4].Int()
	)
	b.width = width
	b.height = height
	b.kLines = kLines
	b.bounceBounds = bounceBounds
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nPoints; i++ {
		x := math.Floor(rand.Float64() * width)
		y := math.Floor(rand.Float64() * height)
		d := rand.Float64() * TAU
		m := lerp(rand.Float64(), 0.1, 0.5)
		b.points = append(b.points, newPoint(x, y, d, m, b))
	}
	b.initCalled = true
	return nil
}
func (b *Board) KLines(this js.Value, args []js.Value) interface{} {
	b.kLines = args[0].Int()
	return nil
}
func (b *Board) Threshold(this js.Value, args []js.Value) interface{} {
	b.threshold = args[0].Float()
	return nil
}

func (b *Board) Points(this js.Value, args []js.Value) interface{} {
	if !b.initCalled {
		return js.ValueOf("board not initialize")
	}
	points := make([]interface{}, len(b.points))
	for i, p := range b.points {
		points[i] = p.Serialize()
	}
	return js.ValueOf(points)
}
func (b *Board) Lines(this js.Value, args []js.Value) interface{} {
	if !b.initCalled {
		return js.ValueOf("board not initialize")
	}
	var lines []interface{}
	for _, p := range b.points {
		for _, l := range p.lines {
			lines = append(lines, l.Serialize())
		}
	}
	return js.ValueOf(lines)
}

func (b *Board) Update(this js.Value, args []js.Value) interface{} {
	if !b.initCalled {
		return js.ValueOf("board not initialized")
	}
	var wg sync.WaitGroup
	for _, p := range b.points {
		wg.Add(1)
		go func(g *sync.WaitGroup, _p *point) {
			defer g.Done()
			_p.update()
		}(&wg, p)
	}
	wg.Wait()
	if b.kLines > 0 {
		for _, p := range b.points {
			wg.Add(1)
			go func(g *sync.WaitGroup, _p *point) {
				defer g.Done()
				_p.updateLines()
			}(&wg, p)
		}
		wg.Wait()
	}
	return nil
}
