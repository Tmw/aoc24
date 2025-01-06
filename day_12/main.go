package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"maps"
	"os"
	"slices"
	"time"
)

var (
	NeighbourNorth = Vector{X: 0, Y: -1}
	NeighbourEast  = Vector{X: 1, Y: 0}
	NeighbourSouth = Vector{X: 0, Y: 1}
	NeighbourWest  = Vector{X: -1, Y: 0}

	OrthogonalNeighbouringDirections = []Vector{
		NeighbourNorth,
		NeighbourEast,
		NeighbourSouth,
		NeighbourWest,
	}

	NeighbourNorthEast = NeighbourNorth.Add(NeighbourEast)
	NeighbourSouthEast = NeighbourSouth.Add(NeighbourEast)
	NeighbourSouthWest = NeighbourNorth.Add(NeighbourWest)
	NeighbourNorthWest = NeighbourSouth.Add(NeighbourWest)
)

type Grid struct {
	Tiles  []byte
	Width  int
	Height int
}
type Borders uint8

func (b Borders) HasBorder(border Borders) bool {
	return b&border > 0
}

const (
	BordersNorth = Borders(1 << iota)
	BordersEast
	BordersSouth
	BordersWest
)

func (g *Grid) BordersAt(loc Vector) Borders {
	var b Borders
	selfTyp := g.At(loc)

	for _, dir := range OrthogonalNeighbouringDirections {
		neighbourLoc := loc.Add(dir)
		if !g.WithinBounds(neighbourLoc) || g.At(neighbourLoc) != selfTyp {
			switch dir {
			case NeighbourNorth:
				b ^= BordersNorth

			case NeighbourEast:
				b ^= BordersEast

			case NeighbourSouth:
				b ^= BordersSouth

			case NeighbourWest:
				b ^= BordersWest
			}
		}
	}

	return b
}

func (g *Grid) WithinBounds(pos Vector) bool {
	return pos.X >= 0 && pos.Y >= 0 &&
		pos.X < g.Width && pos.Y < g.Height
}

func (g *Grid) At(pos Vector) byte {
	return g.Tiles[pos.Y*g.Width+pos.X]
}

func (g *Grid) Print() {
	for idx, tile := range g.Tiles {
		if idx > 0 && idx%g.Width == 0 {
			fmt.Println()
		}

		fmt.Printf("%s", string(tile))
	}

	fmt.Println()
}

func (g *Grid) NeighboursOfType(pos Vector, typ byte) []Vector {
	var neighbours []Vector
	for _, dir := range OrthogonalNeighbouringDirections {
		newPos := pos.Add(dir)
		if g.WithinBounds(newPos) && g.At(newPos) == typ {
			neighbours = append(neighbours, newPos)
		}
	}

	return neighbours
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

type Cluster []Vector

func (c Cluster) Area() int {
	return len(c)
}

func (c Cluster) Sides(g Grid) int {
	var total int

	for _, coord := range c {
		selfTyp := g.At(coord)
		borders := g.BordersAt(coord)
		// fmt.Printf("- checking coord: %+v: borders: %b\n", coord, borders)

		// if we have two connected borders, we have a corner.
		// this only detects corners on the outside.

		// test northeast corner
		if borders.HasBorder(BordersNorth) && borders.HasBorder(BordersEast) {
			total++
		}

		// test southeast corner
		if borders.HasBorder(BordersEast) && borders.HasBorder(BordersSouth) {
			total++
		}

		// test southwest corner
		if borders.HasBorder(BordersSouth) && borders.HasBorder(BordersWest) {
			total++
		}

		// test northwest corner
		if borders.HasBorder(BordersWest) && borders.HasBorder(BordersNorth) {
			total++
		}

		// but that's not all, we also need to detect corners on the inside of the shape,
		// that logic is a bit more convoluted as we'll need to check the neighbours too.

		// detect:
		// EEE <- first E has a corner bottom-right of the first E
		// EXX

		northNeighbour := coord.Add(NeighbourNorth)
		eastNeighbour := coord.Add(NeighbourEast)
		southNeighbour := coord.Add(NeighbourSouth)
		westNeighbour := coord.Add(NeighbourWest)

		northEastNeighbour := coord.Add(NeighbourNorthEast)
		southEastNeighbour := coord.Add(NeighbourSouthEast)
		southWestNeighbour := coord.Add(NeighbourSouthWest)
		northWestNeighbour := coord.Add(NeighbourNorthWest)

		if g.WithinBounds(southNeighbour) && g.At(southNeighbour) == selfTyp &&
			g.WithinBounds(eastNeighbour) && g.At(eastNeighbour) == selfTyp &&
			g.WithinBounds(southEastNeighbour) && g.At(southEastNeighbour) != selfTyp {
			total++
		}

		// detect:
		// EXX
		// EEE <- first E has a corner top-right of the first E
		if g.WithinBounds(northNeighbour) && g.At(northNeighbour) == selfTyp &&
			g.WithinBounds(eastNeighbour) && g.At(eastNeighbour) == selfTyp &&
			g.WithinBounds(northEastNeighbour) && g.At(northEastNeighbour) != selfTyp {
			total++
		}

		// detect:
		// EEE <- last E has a corner bottom-left of the last E
		// XXE
		if g.WithinBounds(southNeighbour) && g.At(southNeighbour) == selfTyp &&
			g.WithinBounds(westNeighbour) && g.At(westNeighbour) == selfTyp &&
			g.WithinBounds(southWestNeighbour) && g.At(southWestNeighbour) != selfTyp {
			total++
		}

		// detect:
		// XXE
		// EEE <- last E has a corner top-left of the last E
		if g.WithinBounds(northNeighbour) && g.At(northNeighbour) == selfTyp &&
			g.WithinBounds(westNeighbour) && g.At(westNeighbour) == selfTyp &&
			g.WithinBounds(northWestNeighbour) && g.At(northWestNeighbour) != selfTyp {
			total++
		}
	}

	return total
}

func (c Cluster) Perimeter(g Grid) int {
	var perimeter int
	for _, loc := range c {
		selfTyp := g.At(loc)
		perimeter += 4 - len(g.NeighboursOfType(loc, selfTyp))
	}

	return perimeter
}

type Clusterer struct {
	grid      Grid
	available map[Vector]struct{}
}

func (c *Clusterer) Init(g Grid) {
	available := map[Vector]struct{}{}
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			available[Vector{X: x, Y: y}] = struct{}{}
		}
	}

	c.grid = g
	c.available = available
}

func (c *Clusterer) nextAvailableStartingLocation() (Vector, bool) {
	next, stop := iter.Pull(maps.Keys(c.available))
	defer stop()
	loc, ok := next()
	if !ok {
		return Vector{}, false
	}

	return loc, ok
}

func (c *Clusterer) MarkUnavailable(locs ...Vector) {
	for _, loc := range locs {
		delete(c.available, loc)
	}
}

func (c *Clusterer) Clusters() []Cluster {
	var clusters []Cluster

	for {
		start, found := c.nextAvailableStartingLocation()
		if !found {
			break
		}

		typ := c.grid.At(start)

		cluster := map[Vector]struct{}{}
		c.findCluster(start, typ, cluster)

		locs := slices.Collect(maps.Keys(cluster))
		c.MarkUnavailable(locs...)
		clusters = append(clusters, locs)
	}

	return clusters
}

func (c *Clusterer) findCluster(start Vector, typ byte, cluster map[Vector]struct{}) {
	cluster[start] = struct{}{}

	neighbours := []Vector{}
	for _, n := range c.grid.NeighboursOfType(start, typ) {
		if _, present := cluster[n]; present {
			continue
		}

		neighbours = append(neighbours, n)
	}

	if len(neighbours) == 0 {
		return
	}

	for _, n := range neighbours {
		if _, present := cluster[n]; present {
			continue
		}

		c.findCluster(n, typ, cluster)
	}
}

func main() {
	grid := parseInput(os.Stdin)

	c := Clusterer{}
	c.Init(grid)
	clusters := c.Clusters()

	start := time.Now()
	fmt.Println("answer part one =", partOne(grid, clusters))
	fmt.Printf("part one took %+v\n", time.Since(start))

	start = time.Now()
	fmt.Println("answer part two =", partTwo(grid, clusters))
	fmt.Printf("part two took %+v", time.Since(start))
}

func partOne(grid Grid, clusters []Cluster) int {
	total := 0
	for _, c := range clusters {
		total += c.Area() * c.Perimeter(grid)
	}

	return total
}

func partTwo(grid Grid, clusters []Cluster) int {
	total := 0
	for _, c := range clusters {
		typ := grid.At(c[0])
		fmt.Printf("======[ cluster %s ]========\n", string(typ))
		sides := c.Sides(grid)
		fmt.Printf(" - cluster %s => sides: %d\n", string(typ), sides)
		total += c.Area() * sides
	}

	return total
}

func parseInput(input io.Reader) Grid {
	var (
		scanner = bufio.NewScanner(input)
		grid    = Grid{}
		x, y    int
	)

	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		tile := scanner.Text()[0]
		if tile == '\n' {
			y++
			if grid.Width == 0 {
				grid.Width = x
			}
			continue
		}

		grid.Tiles = append(grid.Tiles, tile)
		x++
	}

	grid.Height = y
	return grid
}
