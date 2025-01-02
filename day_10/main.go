package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func main() {
	fmt.Println()
}

type Point struct {
	X, Y int
}

type Map struct {
	Cells  []uint8
	Width  int
	Height int

	StartLocations  []Point
	FinishLocations []Point
}

func parseInput(input io.Reader) Map {
	var (
		scanner = bufio.NewScanner(input)
		output  Map
	)

	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		char := scanner.Text()
		if char == "\n" {
			continue
		}

		num, err := strconv.ParseUint(char, 10, 8)
		if err != nil {
			panic(fmt.Errorf("error converting %s to int: %w", char, err))
		}
	}

	return output
}
