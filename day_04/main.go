package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"strings"
)

type Grid struct {
	contents []byte
	width    int
	height   int
}

func (g Grid) String() string {
	return string(g.contents)
}

func (g *Grid) CharAt(x, y int) byte {
	idx := y*g.width + x
	return g.contents[idx]
}

func (g *Grid) Horizontals() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for y := 0; y < g.height; y++ {
			start := y * g.width
			end := start + g.width
			if !yield(g.contents[start:end]) {
				break
			}
		}
	}
}

func (g *Grid) Verticals() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for x := 0; x < g.width; x++ {
			v := make([]byte, 0, g.height)
			for y := 0; y < g.height; y++ {
				v = append(v, g.CharAt(x, y))
			}

			if !yield(v) {
				break
			}
		}
	}
}

// Diagonal returns a single diagonal that starts at the given x and y position.
// the direction integer determines the direction of the slice,
//
// passing +1 iterates through left-to-right
// passing -1 iterates through right-to-left
func (g *Grid) Diagonal(x, y int, direction int) []byte {
	line := []byte{}
	for x >= 0 && y >= 0 && x < g.width && y < g.height {
		line = append(line, g.CharAt(x, y))
		x += direction
		y += 1
	}
	return line
}

func (g *Grid) Diagonals() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		// iterate left-to-right
		x, y := 0, g.height-1
		for y >= 0 && x < g.width {
			if !yield(g.Diagonal(x, y, +1)) {
				return
			}

			if y > 0 {
				y--
			} else {
				x++
			}
		}

		// iterate right-to-left
		x, y = g.width-1, g.height-1
		for y >= 0 && x >= 0 {
			if !yield(g.Diagonal(x, y, -1)) {
				return
			}

			if y > 0 {
				y--
			} else {
				x--
			}
		}
	}
}

func (g *Grid) Subgrid(x, y int, width, height int) Grid {
	var (
		contents = make([]byte, 0, width*height)
	)

	for innerY := range height {
		for innerX := range width {
			contents = append(contents, g.CharAt(x+innerX, y+innerY))
		}
	}

	return Grid{
		contents: contents,
		width:    width,
		height:   height,
	}
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("unable to read input: %v", err)
	}

	grid := ParseGrid(string(input))
	solvePartOne(grid)
	solvePartTwo(grid)

}

func solvePartOne(g Grid) {
	combined := CombineIter(
		g.Horizontals(),
		g.Verticals(),
		g.Diagonals(),
	)

	xmasCount := 0
	for line := range combined {
		xmasCount += strings.Count(string(line), "XMAS")
		xmasCount += strings.Count(string(line), "SAMX")
	}

	fmt.Println("Answer part one:", xmasCount)
}

func solvePartTwo(g Grid) {
	xmasCount := 0
	windowSize := 3

	for offsetY := range g.height - (windowSize - 1) {
		for offsetX := range g.width - (windowSize - 1) {
			sg := g.Subgrid(offsetX, offsetY, windowSize, windowSize)

			diagonalOne := string(sg.Diagonal(0, 0, +1))
			diagonalTwo := string(sg.Diagonal(2, 0, -1))

			if (diagonalOne == "MAS" || diagonalOne == "SAM") &&
				(diagonalTwo == "MAS" || diagonalTwo == "SAM") {
				xmasCount++
			}
		}
	}

	fmt.Println("Answer part two:", xmasCount) // expect 9
}

// CombineIter takes N iter.Seq[T] and returns a single one, concatening
// the passed in iterators into a single one.
func CombineIter[T any, I iter.Seq[T]](i ...I) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, itr := range i {
			for item := range itr {
				if !yield(item) {
					break
				}
			}
		}

	}
}

func ParseGrid(input string) Grid {
	var (
		numLines = strings.Count(input, "\n")
		numCols  = strings.Index(input, "\n")
	)

	return Grid{
		contents: []byte(strings.ReplaceAll(input, "\n", "")),
		width:    numCols,
		height:   numLines,
	}
}
