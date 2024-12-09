package main

import (
	"fmt"
	"io"
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

func valid(r Report) bool {
	d := deltas(r)
	if !allIncrease(d) && !allDecrease(d) {
		return false
	}

	var (
		maxDelta     = slices.Max(d)
		minDelta     = slices.Min(d)
		safeMaxDelta = abs(maxDelta) <= 3 && abs(maxDelta) >= 1
		safeMinDelta = abs(minDelta) >= 1 && abs(minDelta) <= 3
	)

	return safeMaxDelta && safeMinDelta
}

func (r Report) IsSafe(dampenerEnabled bool) bool {
	if valid(r) {
		return true
	}

	if !dampenerEnabled {
		return false
	}

	for p := range slices.Values(permutations(r)) {
		if valid(p) {
			return true
		}
	}

	return false
}

func permutations(r Report) []Report {
	permutated := make([]Report, 0, len(r))

	for idx := range r {
		permutated = append(permutated, without(r, idx))
	}

	return permutated
}

func without(input Report, idx int) Report {
	return slices.Concat(input[0:idx], input[idx+1:])
}

func abs(a int) int {
	if a < 0 {
		return a * -1
	}

	return a
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
		if report.IsSafe(false) {
			safeReports += 1
		}
	}

	fmt.Println("output part one: ", safeReports)

	safeReports = 0
	for report := range slices.Values(reports) {
		if report.IsSafe(true) {
			safeReports += 1
		}
	}

	fmt.Println("output part two: ", safeReports)
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
