package main

import (
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

const (
	TAU = math.Pi * 2
)

type point struct {
	x     float64
	y     float64
	d     float64
	m     float64
	board *Board
}

func (p *point) update() {
	magnitude := time.Since(p.board.lastFrame).Seconds() * p.m * 100
	// magnitude := p.m
	if p.board.bounceBounds {
		// reflect off the edges of the board
		if p.x > p.board.width || p.x < 0 {
			// reflect off wall
			p.d = math.Pi - p.d
		} else if p.y > p.board.height || p.y < 0 {
			// reflect off floor/ceil
			p.d = TAU - p.d
		}
		p.x = p.x + (math.Cos(p.d) * magnitude)
		p.y = p.y + (math.Sin(p.d) * magnitude)
	} else {
		// pass through the edges and come out the other side
		var xCorrection, yCorrection float64
		if p.x > p.board.width {
			xCorrection = -p.board.width
		} else if p.x < 0 {
			xCorrection = p.board.width
		}
		if p.y > p.board.height {
			yCorrection = -p.board.height
		} else if p.y < 0 {
			yCorrection = p.board.height
		}

		p.x = p.x + (math.Cos(p.d) * magnitude) + xCorrection
		p.y = p.y + (math.Sin(p.d) * magnitude) + yCorrection
	}
}

type Board struct {
	initCalled   bool
	width        float64
	height       float64
	points       []*point
	bounceBounds bool
	lastFrame    time.Time
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
	)
	b.width = width
	b.height = height
	b.bounceBounds = bounceBounds
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nPoints; i++ {
		b.points = append(b.points, &point{
			x:     math.Floor(rand.Float64() * width),
			y:     math.Floor(rand.Float64() * height),
			d:     rand.Float64() * TAU,
			m:     rand.Float64(),
			board: b,
		})
	}
	b.initCalled = true
	return nil
}

func (b *Board) Points(this js.Value, args []js.Value) interface{} {
	if !b.initCalled {
		return js.ValueOf("board not initialize")
	}
	b.lastFrame = time.Now()
	points := make([]interface{}, len(b.points))
	for i, p := range b.points {
		points[i] = map[string]interface{}{
			"x": p.x,
			"y": p.y,
		}
	}
	return js.ValueOf(points)
}

func (b *Board) Update(this js.Value, args []js.Value) interface{} {
	if !b.initCalled {
		return js.ValueOf("board not initialized")
	}
	for _, p := range b.points {
		p.update()
	}
	return nil
}

func main() {
	c := make(chan struct{}, 0)
	b := &Board{}
	functions := map[string]js.Func{
		"init":         js.FuncOf(b.Init),
		"getPoints":    js.FuncOf(b.Points),
		"updatePoints": js.FuncOf(b.Update),
	}
	for name, f := range functions {
		js.Global().Set(name, f)
	}
	<-c
}
