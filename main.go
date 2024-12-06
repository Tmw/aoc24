package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")
	separator := "   "
	listA := make([]int, 0, len(lines))
	listB := make([]int, 0, len(lines))

	for line := range slices.Values(lines) {
		parts := strings.Split(line, separator)
		if len(parts) != 2 {
			continue
		}

		itemA, _ := strconv.Atoi(parts[0])
		itemB, _ := strconv.Atoi(parts[1])

		listA = append(listA, itemA)
		listB = append(listB, itemB)
	}

	slices.Sort(listA)
	slices.Sort(listB)
	diff := float64(0)

	freqMap := make(map[int]int)
	for _, num := range listB {
		if _, ok := freqMap[num]; !ok {
			freqMap[num] = 0
		}

		freqMap[num] += 1
	}

	for idx := range listA {
		left, right := listA[idx], listB[idx]
		diff += math.Abs(float64(left - right))
	}

	fmt.Printf("output Part One: %d\n", int(diff))

	partB := 0
	for _, num := range listA {
		if multiplier, found := freqMap[num]; found {
			partB += num * multiplier
		}
	}
	fmt.Printf("output Part Two: %d\n", partB)
}
