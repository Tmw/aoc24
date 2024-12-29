package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type Vector struct {
	X, Y int
}

func (v Vector) Sub(v2 Vector) Vector {
	return Vector{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

type Arena struct {
	width     int
	height    int
	locByFreq map[rune][]Vector
}

func NewArena() Arena {
	return Arena{
		locByFreq: make(map[rune][]Vector),
	}
}

func (a *Arena) withinBounds(v Vector) bool {
	return v.X >= 0 && v.X < a.width &&
		v.Y >= 0 && v.Y < a.height
}

func main() {
	arena := parseInput(os.Stdin)

	start := time.Now()
	fmt.Println("answer part one =", solve(arena, false))
	fmt.Printf("part one took %+v\n", time.Since(start))

	start = time.Now()
	fmt.Println("answer part two =", solve(arena, true))
	fmt.Printf("part two took %+v", time.Since(start))
}

func solve(a Arena, resonance bool) int {
	antiNodes := map[Vector]struct{}{}
	for _, locs := range a.locByFreq {
		for _, pair := range pairs(locs) {

			// When we take resonance into account, the antennas themselves
			// become antinodes too.
			if resonance {
				antiNodes[pair[0]] = struct{}{}
				antiNodes[pair[1]] = struct{}{}
			}

			var (
				d1        = pair[0].Sub(pair[1])
				d2        = pair[1].Sub(pair[0])
				antiNode1 = pair[0]
				antiNode2 = pair[1]
			)

			for {
				antiNode1 = antiNode1.Add(d1)
				antiNode2 = antiNode2.Add(d2)
				shouldContinue := false

				if a.withinBounds(antiNode1) {
					antiNodes[antiNode1] = struct{}{}
					shouldContinue = true
				}

				if a.withinBounds(antiNode2) {
					antiNodes[antiNode2] = struct{}{}
					shouldContinue = true
				}

				// when we take resonance into account, the antinodes
				// will be present at each interval and won't stop afer the first occurance.
				if !resonance || !shouldContinue {
					break
				}
			}
		}
	}
	return len(antiNodes)
}

func pairs[T any](s []T) [][2]T {
	size := len(s) * (len(s) - 1) / 2
	res := make([][2]T, 0, size)

	for a := 0; a <= len(s)-1; a++ {
		for b := a + 1; b < len(s); b++ {
			res = append(res, [2]T{s[a], s[b]})
		}
	}

	return res
}

func parseInput(input io.Reader) Arena {
	var (
		arena    = NewArena()
		scanner  = bufio.NewScanner(input)
		row, col int
	)

	for scanner.Scan() {
		cells := scanner.Text()
		arena.width = len(cells)
		col = 0

		for _, cell := range cells {
			if cell == '.' {
				col++
				continue
			}

			if _, present := arena.locByFreq[cell]; !present {
				arena.locByFreq[cell] = []Vector{}
			}

			tmp := arena.locByFreq[cell]
			tmp = append(tmp, Vector{X: col, Y: row})
			arena.locByFreq[cell] = tmp
			col++
		}
		row++
	}

	arena.height = row
	return arena
}
