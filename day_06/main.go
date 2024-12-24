package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"time"
)

const (
	MaxIterationPartOne = 10_000
	MaxIterationPartTwo = 20_000
)

func main() {
	grid, startingPos := parseInput(os.Stdin)

	start := time.Now()

	path, numUniqueLocs := partOne(grid, startingPos)
	fmt.Println("Answer part one = ", numUniqueLocs)
	fmt.Println("Answer part two = ", partTwo(grid, startingPos, path))
}

func parseInput(input io.Reader) (Grid, Vector) {
	var (
		res         Grid
		pos         Vector
		startingPos Vector
	)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		res.width = len(line)
		for _, char := range line {
			switch char {
			case '.':
				res.cells = append(res.cells, CellTypeOpen)

			case '#':
				res.cells = append(res.cells, CellTypeBlocked)

			case '^':
				res.cells = append(res.cells, CellTypeOpen)
				startingPos = pos
			}

			pos.X++
		}

		pos.Y++
		pos.X = 0
	}
	res.height = pos.Y

	return res, startingPos
}

func partOne(grid Grid, startingPos Vector) ([]Vector, int) {
	w := Walker{
		grid:      grid,
		pos:       startingPos,
		heading:   HeadingNorth,
		stepsLeft: MaxIterationPartOne,
	}

	_, path := w.WalkToCompletion()
	return path, len(unique(path))
}

func unique[T comparable](s []T) []T {
	unique := map[T]struct{}{}
	for _, loc := range s {
		unique[loc] = struct{}{}
	}
	return slices.Collect(maps.Keys(unique))
}

func partTwo(grid Grid, startingPos Vector, path []Vector) int {
	candidates := []Vector{}

	for _, coord := range slices.Backward(path) {
		g := grid.Clone()
		g.SetCellAt(coord, CellTypeBlocked)
		w := Walker{
			grid:      g,
			pos:       startingPos,
			heading:   HeadingNorth,
			stepsLeft: MaxIterationPartTwo,
		}

		condition, _ := w.WalkToCompletion()
		if condition == StepResultsExhaustedMaxSteps {
			// infinite loop detected,..
			candidates = append(candidates, coord)
		}
	}

	return len(unique(candidates))
}
