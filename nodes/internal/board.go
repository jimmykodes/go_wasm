package internal

import (
	"math"
	"math/rand"
	"sync"
	"syscall/js"
	"time"
)

type Board struct {
	width        float64
	height       float64
	kLines       int
	points       []*point
	wells        []*Well
	bounceBounds bool
	threshold    float64
	nPoints      int
}

func NewBoard(config js.Value) *Board {
	var (
		width        = config.Get("width")
		height       = config.Get("height")
		nPoints      = config.Get("nPoints")
		bounceBounds = config.Get("bounceBounds")
		kLines       = config.Get("kLines")
		threshold    = config.Get("threshold")
	)
	if width.IsUndefined() {
		width = js.ValueOf(400)
	}
	if height.IsUndefined() {
		width = js.ValueOf(400)
	}
	if nPoints.IsUndefined() {
		nPoints = js.ValueOf(100)
	}
	if kLines.IsUndefined() {
		kLines = js.ValueOf(2)
	}
	if bounceBounds.IsUndefined() {
		bounceBounds = js.ValueOf(false)
	}
	if threshold.IsUndefined() {
		threshold = js.ValueOf(20)
	}
	b := &Board{
		width:        width.Float(),
		height:       height.Float(),
		nPoints:      nPoints.Int(),
		kLines:       kLines.Int(),
		bounceBounds: bounceBounds.Bool(),
		threshold:    threshold.Float() * 1000,
		wells: []*Well{
			{
				x:        width.Float() / 2,
				y:        height.Float() / 2,
				strength: 10,
			},
		},
	}
	return b
}

func (b *Board) Serializer() interface{} {
	return map[string]interface{}{
		"initPoints": js.FuncOf(b.InitPoints),
		"kLines":     js.FuncOf(b.KLines),
		"getPoints":  js.FuncOf(b.Points),
		"getLines":   js.FuncOf(b.Lines),
		"update":     js.FuncOf(b.Update),
	}
}

func (b *Board) InitPoints(this js.Value, args []js.Value) interface{} {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.nPoints; i++ {
		x := math.Floor(rand.Float64() * b.width)
		y := math.Floor(rand.Float64() * b.height)
		d := rand.Float64() * TAU
		m := lerp(rand.Float64(), 0.1, 0.5)
		b.points = append(b.points, newPoint(x, y, d, m, b))
	}
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
	points := make([]interface{}, len(b.points))
	for i, p := range b.points {
		points[i] = p.Serialize()
	}
	return js.ValueOf(points)
}
func (b *Board) Lines(this js.Value, args []js.Value) interface{} {
	var lines []interface{}
	for _, p := range b.points {
		for _, l := range p.lines {
			lines = append(lines, l.Serialize())
		}
	}
	return js.ValueOf(lines)
}

func (b *Board) Update(this js.Value, args []js.Value) interface{} {
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
