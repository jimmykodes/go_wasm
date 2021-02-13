package internal

type line struct {
	start     *point
	end       *point
	squareMag float64
}

func (l line) Serialize() interface{} {
	return map[string]interface{}{
		"start": map[string]interface{}{
			"x": l.start.x,
			"y": l.start.y,
		},
		"end": map[string]interface{}{
			"x": l.end.x,
			"y": l.end.y,
		},
	}
}
