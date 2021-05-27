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
	p.updateMag()
	if p.board.bounceBounds {
		// reflect off the edges of the board
		if p.x > p.board.width || p.x < 0 {
			// reflect off wall
			p.d = math.Pi - p.d
		} else if p.y > p.board.height || p.y < 0 {
			// reflect off floor/ceil
			p.d = TAU - p.d
		}
		p.x = p.x + (math.Cos(p.d) * p.m)
		p.y = p.y + (math.Sin(p.d) * p.m)
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

		p.x = p.x + (math.Cos(p.d) * p.m) + xCorrection
		p.y = p.y + (math.Sin(p.d) * p.m) + yCorrection
	}
}

// todo: figure out why this function doesn't work.
func (p *point) updateMag() {
	for _, well := range p.board.wells {
		sqrMag := squareMag(p, well)
		dy := p.y - well.y
		dx := p.x - well.x
		var direction float64
		if dx != 0 {
			direction = math.Atan(dy/dx)
		} else {
			if p.y > well.y {
				direction = math.Pi * 1.5
			} else {
				direction = math.Pi * .5
			}
		}
		if p.x > well.x {
			if p.y > well.y {
				direction = math.Pi + direction
			} else {
				direction = math.Pi - direction
			}
		} else {
			if p.y > well.y {
				direction = TAU - direction
			}
		}
		strength := well.strength / sqrMag
		newM := m3(p.m, strength, p.d, direction)
		p.d = t3(p.m, strength, newM, p.d, direction)
		p.m = newM
	}
}

func (p *point) updateLines() {
	p.lines = nil
	for _, point := range p.board.points {
		if point.uuid == p.uuid {
			// don't create a line to ourself
			continue
		}
		if sm := squareMag(p, point); sm < p.board.threshold || p.board.threshold < 0 {
			// don't create entries in the lines array we are inevitably going to filter out
			// threshold value of < 0 means infinite threshold
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
func (p *point) X() float64 {
	return p.x
}
func (p *point) Y() float64 {
	return p.y
}
func (p *point) Mag() float64 {
	return p.m
}
func (p *point) Serialize() interface{} {
	return map[string]interface{}{
		"x": p.x,
		"y": p.y,
	}
}
