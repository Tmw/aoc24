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

type Robot struct {
	Pos Vector
	Vel Vector
}

type Quadrant int

const (
	QuadrantNone = Quadrant(iota)
	QuadrantTopLeft
	QuadrantTopRight
	QuadrantBottomLeft
	QuadrantBottomRight
)

const (
	mapWidth  = 101 // 11 for example input
	mapHeight = 103 // 7 for example input
	iters     = 100
)

func determineQuadrant(x, y int) Quadrant {
	switch {
	case x < mapWidth/2 && y < mapHeight/2:
		return QuadrantTopLeft

	case x > mapWidth/2 && y < mapHeight/2:
		return QuadrantTopRight

	case x < mapWidth/2 && y > mapHeight/2:
		return QuadrantBottomLeft

	case x > mapWidth/2 && y > mapHeight/2:
		return QuadrantBottomRight
	}

	// if we couldn't match any of the above statements,
	// the robot is on the border between quadrants.
	// these do not count.
	return QuadrantNone
}

func main() {
	robots := parseInput(os.Stdin)

	start := time.Now()
	fmt.Println("answer part one =", partOne(robots))
	fmt.Printf("part one took %+v\n", time.Since(start))

	// start = time.Now()
	// fmt.Println("answer part two =", cost(robots))
	// fmt.Printf("part two took %+v", time.Since(start))
}

func partOne(robots []Robot) int {
	countPerQuadrant := map[Quadrant]int{}

	for _, r := range robots {
		// find the final X and Y after N iterations in one go using modulo
		x := ((r.Pos.X+r.Vel.X*iters)%mapWidth + mapWidth) % mapWidth
		y := ((r.Pos.Y+r.Vel.Y*iters)%mapHeight + mapHeight) % mapHeight

		// find to which quadrant the robot belongs based on its final X and Y
		// coordinate and increment the counter that belongs to the quadrant.
		q := determineQuadrant(x, y)
		countPerQuadrant[q] += 1
	}

	total := 1
	for q, v := range countPerQuadrant {
		if q == QuadrantNone {
			continue
		}
		total *= v
	}

	return total
}

func parseInput(input io.Reader) []Robot {
	var (
		scanner = bufio.NewScanner(input)
		robots  = []Robot{}
	)

	for scanner.Scan() {
		line := scanner.Text()

		var m Robot
		if line == "" {
			continue
		}

		_, err := fmt.Sscanf(line, "p=%d,%d v=%d,%d", &m.Pos.X, &m.Pos.Y, &m.Vel.X, &m.Vel.Y)
		if err != nil {
			panic(fmt.Errorf("unable to parse line '%s': %w", line, err))
		}

		robots = append(robots, m)
	}

	return robots
}
