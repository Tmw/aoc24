package main

type Walker struct {
	grid      Grid
	pos       Vector
	heading   Heading
	stepsLeft int
}

type StepResult int

const (
	StepResultOutOfBounds StepResult = iota
	StepResultsNonWalkable
	StepResultsExhaustedMaxSteps
	StepResultOK
)

func (w *Walker) Step() StepResult {
	nextPos := w.NextPos()

	if !w.grid.WithinBounds(nextPos) {
		return StepResultOutOfBounds
	}

	if !w.grid.CellAtPositionWalkable(nextPos) {
		return StepResultsNonWalkable
	}

	if w.stepsLeft <= 0 {
		return StepResultsExhaustedMaxSteps
	}

	w.stepsLeft -= 1
	w.MoveTo(nextPos)
	return StepResultOK
}

func (w *Walker) NextPos() Vector {
	return w.pos.Add(w.heading.Vector())
}

func (w *Walker) MoveTo(pos Vector) {
	w.pos = pos
}

// walk until one of the results is reached
func (w *Walker) WalkToCompletion() (StepResult, []Vector) {
	var path []Vector

	for {
		path = append(path, w.pos)
		res := w.Step()
		if res == StepResultOutOfBounds {
			return res, path
		}

		if res == StepResultsExhaustedMaxSteps {
			return res, path
		}

		if res == StepResultsNonWalkable {
			w.heading = w.heading.RotateClockwise()
		}
	}
}
