package main

import "iter"

type Walker struct {
	grid    Grid
	pos     Vector
	heading Heading
}

type Step struct {
	Pos     Vector
	Heading Heading
}

func (w *Walker) Walk() iter.Seq[Step] {
	return func(yield func(Step) bool) {
		for {
			s := Step{
				Pos:     w.pos,
				Heading: w.heading,
			}

			if !yield(s) {
				return
			}

			nextPos := w.nextPos()
			if !w.grid.WithinBounds(nextPos) {
				return
			}

			if !w.grid.CellAtPositionWalkable(nextPos) {
				w.heading = w.heading.RotateClockwise()
				continue
			}

			w.moveTo(nextPos)
		}
	}
}

func (w *Walker) nextPos() Vector {
	return w.pos.Add(w.heading.Vector())
}

func (w *Walker) moveTo(pos Vector) {
	w.pos = pos
}
