package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	CellTrailHead = uint8(0)
	CellTrailPeak = uint8(9)

	DirectionNorth = Vector{X: 0, Y: -1}
	DirectionEast  = Vector{X: 1, Y: 0}
	DirectionSouth = Vector{X: 0, Y: 1}
	DirectionWest  = Vector{X: -1, Y: 0}

	NeighbouringDirections = []Vector{
		DirectionNorth,
		DirectionEast,
		DirectionSouth,
		DirectionWest,
	}
)

func main() {
	m := parseInput(os.Stdin)

	start := time.Now()
	fmt.Println(" answer part one =", partOne(m))
	fmt.Printf("part one took %+v\n", time.Since(start))

	start = time.Now()
	fmt.Println("answer part two =", partTwo(m))
	fmt.Printf("part two took %+v", time.Since(start))
}

func partOne(m Map) int {
	sum := 0
	for _, h := range m.TrailHeads {
		p := ReachablePeaks(m, h, []Vector{})
		if len(p) > 0 {
			sum += len(unique(p))
		}
	}

	return sum
}

// for part two we want to know how many distinct paths we can find to a top.
// using my solution of part one already returns the peak multiple times if it
// would be reachable multiple times, so instead of running the results through unique,
// we'll just make a frequency map and add the totals.
func partTwo(m Map) int {
	sum := 0
	for _, h := range m.TrailHeads {
		p := ReachablePeaks(m, h, []Vector{})
		for n := range maps.Values(freq(p)) {
			sum += n
		}
	}

	return sum
}

func freq[T comparable](s []T) map[T]int {
	res := make(map[T]int)

	for _, item := range s {
		if _, found := res[item]; !found {
			res[item] = 0
		}

		res[item]++
	}

	return res
}

func unique[T comparable](s []T) []T {
	var res []T
	known := make(map[T]struct{})
	for _, item := range s {
		if _, found := known[item]; found {
			continue
		}

		known[item] = struct{}{}
		res = append(res, item)
	}

	return res
}

func ReachablePeaks(m Map, pos Vector, peaks []Vector) []Vector {
	if m.CellAtPos(pos) == CellTrailPeak {
		return append(peaks, pos)
	}

	var reachable []Vector

	for _, n := range m.ValidNeighbours(pos) {
		if p := ReachablePeaks(m, n, peaks); len(p) > 0 {
			reachable = append(reachable, p...)
		}
	}

	return reachable
}

type Vector struct {
	X, Y int
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

type Map struct {
	Cells  []uint8
	Width  int
	Height int

	TrailHeads []Vector
}

func (m *Map) PosInBounds(pos Vector) bool {
	return pos.X >= 0 && pos.Y >= 0 &&
		pos.X < m.Width && pos.Y < m.Height
}

func (m *Map) CellAtPos(pos Vector) uint8 {
	if !m.PosInBounds(pos) {
		panic(fmt.Errorf("position out of bounds"))
	}

	idx := pos.Y*m.Width + pos.X
	return m.Cells[idx]
}

// ValidNeighbours returns the neighbours based off of the current
// position that are both in the map as well as exactly one level higher.
func (m *Map) ValidNeighbours(pos Vector) []Vector {
	var (
		res        = make([]Vector, 0, 4)
		currentVal = m.CellAtPos(pos)
	)

	for dir := range slices.Values(NeighbouringDirections) {
		neighbour := pos.Add(dir)

		// neighbouring cell within bounds and exactly one level higher?
		if m.PosInBounds(neighbour) && m.CellAtPos(neighbour)-currentVal == 1 {
			res = append(res, neighbour)
		}
	}

	return res
}

func (m *Map) String() string {
	var s strings.Builder
	for idx, v := range m.Cells {

		if idx > 0 && idx%m.Width == 0 {
			s.WriteString("\n")
		}
		s.WriteString(fmt.Sprintf("%d", v))
	}

	return s.String()
}

func parseInput(input io.Reader) Map {
	var (
		scanner  = bufio.NewScanner(input)
		output   Map
		col, row int
	)

	for scanner.Scan() {
		cells := scanner.Text()
		output.Width = len(cells)
		col = 0

		for _, cell := range cells {
			num, err := strconv.ParseUint(string(cell), 10, 8)
			if err != nil {
				panic(fmt.Errorf("error converting %s to int: %w", string(cell), err))
			}

			// find trailheads
			if uint8(num) == CellTrailHead {
				output.TrailHeads = append(output.TrailHeads, Vector{X: col, Y: row})
			}

			output.Cells = append(output.Cells, uint8(num))
			col++
		}
		row++
	}

	output.Height = row
	return output
}
