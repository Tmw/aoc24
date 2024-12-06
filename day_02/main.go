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

func deltas(input []int) []int {
	var (
		deltas = make([]int, 0, len(input))
		prev   = 0
	)

	for num := range slices.Values(input) {
		delta := num - prev
		deltas = append(deltas, delta)
		prev = num
	}

	return deltas[1:]
}

type Report []int

func allIncrease(deltas []int) bool {
	for val := range slices.Values(deltas) {
		if val <= 0 {
			return false
		}
	}

	return true
}

func allDecrease(deltas []int) bool {
	for val := range slices.Values(deltas) {
		if val >= 0 {
			return false
		}
	}
	return true
}

func (r Report) IsSafe() bool {
	d := deltas(r)
	if !allIncrease(d) && !allDecrease(d) {
		return false
	}

	var (
		maxDelta     = slices.Max(d)
		minDelta     = slices.Min(d)
		safeMaxDelta = intAbs(maxDelta) <= 3 && intAbs(maxDelta) >= 1
		safeMinDelta = intAbs(minDelta) >= 1 && intAbs(minDelta) <= 3
	)

	return safeMaxDelta && safeMinDelta
}

func intAbs(a int) int {
	return int(math.Abs(float64(a)))
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Errorf("error reading input: %w", err))
	}

	var (
		lines   = strings.Split(string(input), "\n")
		reports = make([]Report, 0, len(lines))
	)

	for line := range slices.Values(lines) {
		if len(line) == 0 {
			continue
		}
		report := parseReport(line)
		reports = append(reports, report)
	}

	safeReports := 0
	for report := range slices.Values(reports) {
		if report.IsSafe() {
			safeReports += 1
		}
	}

	fmt.Println("output part one: ", safeReports)
}

func parseReport(line string) Report {
	var (
		parts  = strings.Split(line, " ")
		report = make(Report, 0, len(parts))
	)

	for part := range slices.Values(parts) {
		val, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		report = append(report, val)
	}

	return report
}
