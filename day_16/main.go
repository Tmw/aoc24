package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"slices"
	"strings"
	"time"
)

type Vector struct {
	X, Y int
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v Vector) Sub(v2 Vector) Vector {
	return Vector{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

type TileType rune

const (
	TileTypeWall  = TileType('#')
	TileTypeOpen  = TileType('.')
	TileTypeStart = TileType('S')
	TileTypeEnd   = TileType('E')
)

var (
	DirectionNorth = Vector{X: 0, Y: -1}
	DirectionEast  = Vector{X: 1, Y: 0}
	DirectionSouth = Vector{X: 0, Y: 1}
	DirectionWest  = Vector{X: -1, Y: 0}
)

type Grid struct {
	width  int
	height int
	tiles  []TileType
}

type Reindeer struct {
	pos Vector
	dir Vector
}

func (g *Grid) Print(start, end, reindeer Vector) {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			pos := Vector{X: x, Y: y}

			if pos == reindeer {
				fmt.Print("@")
				continue
			}

			if pos == start {
				fmt.Print("S")
				continue
			}

			if pos == end {
				fmt.Print("E")
				continue
			}

			tile := g.TileAt(pos)
			fmt.Printf("%s", string(tile))
		}
		fmt.Println()
	}
}

func (g *Grid) PrintWithPath(start, end Vector, path []Vector) {
	res := make([]rune, g.width*g.height)

	// draw "normal" map with start, end and current pos of reindeer?..
	// and all the other tiles as open.
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			pos := Vector{X: x, Y: y}
			idx := y*g.width + x

			if pos == start {
				res[idx] = 'S'
				continue
			}

			if pos == end {
				res[idx] = 'E'
				continue
			}

			tile := g.TileAt(pos)
			res[idx] = rune(tile)
		}
	}

	direction := DirectionEast
	prevPos := start
	for _, pos := range path {
		idx := pos.Y*g.width + pos.X

		if pos != prevPos {
			direction = pos.Sub(prevPos)
		}

		prevPos = pos
		switch direction {
		case DirectionNorth:
			res[idx] = '^'

		case DirectionWest:
			res[idx] = '<'

		case DirectionSouth:
			res[idx] = 'v'

		case DirectionEast:
			res[idx] = '>'
		}
		continue
	}

	var b strings.Builder
	for idx, cell := range res {
		if idx > 0 && idx%g.width == 0 {
			b.WriteRune('\n')
		}

		switch cell {
		case '>', '<', 'v', '^':
			b.WriteString(fmt.Sprintf("\033[32m%s\033[0m", string(cell)))

		default:
			b.WriteString(fmt.Sprintf("\033[90m%s\033[0m", string(cell)))
		}
	}
	fmt.Print(b.String())
}

func (g *Grid) TileAt(pos Vector) TileType {
	idx := g.width*pos.Y + pos.X
	return g.tiles[idx]
}

func (g *Grid) WalkableTilesSurrounding(pos Vector) []Vector {
	var walkable []Vector
	surroundings := []Vector{
		DirectionNorth,
		DirectionEast,
		DirectionSouth,
		DirectionWest,
	}

	for _, dir := range surroundings {
		newPos := pos.Add(dir)
		if g.TileAt(newPos) == TileTypeOpen {
			walkable = append(walkable, newPos)
		}
	}

	return walkable
}

func main() {
	grid, start, end := parseInput(os.Stdin)

	pf := NewPathFinder(PathFinderOpts{
		NeighboursFn: grid.WalkableTilesSurrounding,
		HeuristicFn: func(l Vector) int {
			return manhattan(l, end)
		},
		ReachedFinishFn: func(l Vector) bool {
			return l == end
		},
	})

	path := pf.Path(start, func(c Candidate, path []Vector) {
		fmt.Printf("candidate = %+v\n", c)
		grid.PrintWithPath(start, end, path)
		time.Sleep(200 * time.Millisecond)
	})
	fmt.Println("final path:", costForPath(path))
	grid.PrintWithPath(start, end, path)
}

func manhattan(a, b Vector) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(a int) int {
	return int(math.Abs(float64(a)))
}

type Set[K comparable] struct {
	items   []K
	present map[K]struct{}
}

func NewSet[K comparable]() Set[K] {
	return Set[K]{
		items:   make([]K, 0),
		present: make(map[K]struct{}),
	}
}

func (s *Set[K]) Insert(item K) bool {
	if _, present := s.present[item]; present {
		// did not insert, was already present
		return false
	}

	s.present[item] = struct{}{}
	s.items = append(s.items, item)

	return true
}

func (s *Set[K]) Includes(item K) bool {
	_, present := s.present[item]
	return present
}

func (s *Set[K]) Slice() []K {
	return s.items
}

func (s *Set[K]) Clone() Set[K] {
	return Set[K]{
		items:   slices.Clone(s.items),
		present: maps.Clone(s.present),
	}
}

func costForPath(path []Vector) int {
	dir := DirectionEast
	prev := path[0]
	zeroDir := Vector{}

	// base cost eq. length of the path
	cost := len(path) - 1

	// then apply penalty for each turn
	for _, loc := range path {
		newDir := loc.Sub(prev)
		if dir != newDir && newDir != zeroDir {
			cost += 1000
			dir = newDir
		}
		prev = loc
	}

	return cost
}

func parseInput(input io.Reader) (Grid, Vector, Vector) {
	var (
		scanner = bufio.NewScanner(input)
		grid    Grid
		start   Vector
		end     Vector
	)

	for scanner.Scan() {
		line := scanner.Text()
		grid.width = 0

		for _, c := range line {
			switch TileType(c) {
			case TileTypeWall, TileTypeOpen:
				grid.tiles = append(grid.tiles, TileType(c))

			case TileTypeStart:
				grid.tiles = append(grid.tiles, TileTypeOpen)
				start = Vector{X: grid.width, Y: grid.height}

			case TileTypeEnd:
				grid.tiles = append(grid.tiles, TileTypeOpen)
				end = Vector{X: grid.width, Y: grid.height}
			}

			grid.width++
		}

		grid.height++
	}

	return grid, start, end
}
