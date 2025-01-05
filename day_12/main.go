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

	NeighbouringDirections = []Vector{
		NeighbourNorth,
		NeighbourEast,
		NeighbourSouth,
		NeighbourWest,
	}
)

type Grid struct {
	Tiles  []byte
	Width  int
	Height int
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
	for _, dir := range NeighbouringDirections {
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
	// come up with an algorithm to determine the number of sides of the cluster.

	return 0
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
		total += c.Area() * c.Sides(grid)
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
