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

func partOne(grid Grid, startingPos Vector) ([]Step, int) {
	w := Walker{
		grid:    grid,
		pos:     startingPos,
		heading: HeadingNorth,
	}

	path := slices.Collect(w.Walk())
	coords := make([]Vector, 0, len(path))
	for _, step := range path {
		coords = append(coords, step.Pos)
	}
	return path, len(unique(coords))
}

func unique[T comparable](s []T) []T {
	unique := map[T]struct{}{}
	for _, loc := range s {
		unique[loc] = struct{}{}
	}
	return slices.Collect(maps.Keys(unique))
}

func partTwo(grid Grid, startingPos Vector, path []Step) int {
	var (
		candidateChan = make(chan Vector)
		wg            sync.WaitGroup
	)

	for _, step := range slices.Backward(path) {
		wg.Add(1)
		go func(grid Grid, step Step) {
			defer wg.Done()
			g := grid.Clone()
			g.SetCellAt(step.Pos, CellTypeBlocked)
			w := Walker{
				grid:    g,
				pos:     startingPos,
				heading: HeadingNorth,
			}

			walked := map[Step]struct{}{}
			for woot := range w.Walk() {
				if _, found := walked[woot]; found {
					candidateChan <- step.Pos
					return
				}
				walked[woot] = struct{}{}
			}
		}(grid, step)
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
