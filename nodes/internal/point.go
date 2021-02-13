package internal

import (
	"math"
	"sort"

	"github.com/google/uuid"
)

type point struct {
	uuid  uuid.UUID
	x     float64
	y     float64
	d     float64
	m     float64
	board *Board
	lines []*line
}

func newPoint(x float64, y float64, d float64, m float64, board *Board) *point {
	return &point{x: x, y: y, d: d, m: m, board: board, uuid: uuid.New()}
}

func (p *point) update() {
	magnitude := p.m
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
func (p *point) updateLines() {
	p.lines = nil
	for _, point := range p.board.points {
		if point.uuid == p.uuid {
			// don't create a line to ourself
			continue
		}
		if sm := squareMag(p, point); sm < p.board.threshold {
			// don't create entries in the lines array we are inevitably going to filter out
			p.lines = append(p.lines, &line{
				start:     p,
				end:       point,
				squareMag: sm,
			})
		}
	}
	sort.SliceStable(p.lines, func(i, j int) bool {
		return p.lines[i].squareMag < p.lines[j].squareMag
	})
	if len(p.lines) > p.board.kLines {
		// only take the smallest k lines
		p.lines = p.lines[:p.board.kLines]
	}
}

func (p *point) Serialize() interface{} {
	return map[string]interface{}{
		"x": p.x,
		"y": p.y,
	}
}
