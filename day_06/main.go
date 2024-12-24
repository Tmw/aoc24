package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"sync"
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

	duration := time.Since(start)
	fmt.Println("found answer in", duration)
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
	var (
		candidateChan = make(chan Vector)
		wg            sync.WaitGroup
	)

	// unbounded concurrency, but we cut the solve time in half =)
	for _, coord := range slices.Backward(path) {
		wg.Add(1)
		go func(grid Grid, coord Vector) {
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
				// steps exhausted, likely infinite loop
				candidateChan <- coord
			}

			wg.Done()
		}(grid, coord)
	}

	candidates := []Vector{}
	go func() {
		for coord := range candidateChan {
			candidates = append(candidates, coord)
		}
	}()

	wg.Wait()
	close(candidateChan)
	return len(unique(candidates))
}
