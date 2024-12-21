package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type PrioMap map[int][]int
type Update []int

// Figure out if the update at hand matches the sorting rules
func (u Update) IsCorrectlySorted(rules PrioMap) bool {
	for pos, num := range u {
		expectedPos := abs(lastIndex(u)-lastIndex(rules[num])) - 1
		if pos != expectedPos {
			return false
		}
	}

	return true
}

func (u Update) MiddlePageNumber() int {
	return u[floor(len(u)/2)]
}

func (u Update) FixOrder(rules PrioMap) {
	slices.SortFunc(u, func(a, b int) int {
		return len(rules[a]) - len(rules[b])
	})
}

// generate a new PrioMap, scoped to the numbers in the update itself.
func (p PrioMap) Scope(nums []int) PrioMap {
	m := PrioMap{}

	for _, n := range nums {
		if _, exists := p[n]; !exists {
			continue
		}

		if _, exists := m[n]; !exists {
			m[n] = []int{}
		}

		for _, page := range p[n] {
			if slices.Contains(nums, page) {
				m[n] = append(m[n], page)
			}
		}
	}

	return m
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("unable to read input: %v", err)
	}

	rules, updates, ok := strings.Cut(string(input), "\n\n")
	if !ok {
		log.Fatalf("unable to parse in head and updates section")
	}

	var (
		prioMap       = parseRules(rules)
		parsedUpdates = parseUpdates(updates)
	)

	fmt.Println("Answer part one = ", partOne(parsedUpdates, prioMap))
	fmt.Println("Answer part two = ", partTwo(parsedUpdates, prioMap))
}

func partOne(updates []Update, rules PrioMap) int {
	sum := 0
	for u := range slices.Values(updates) {
		scoped := rules.Scope(u)
		if u.IsCorrectlySorted(scoped) {
			sum += u.MiddlePageNumber()
		}
	}

	return sum
}

func partTwo(updates []Update, rules PrioMap) int {
	sum := 0
	for u := range slices.Values(updates) {
		scoped := rules.Scope(u)
		if !u.IsCorrectlySorted(scoped) {
			u.FixOrder(scoped)
			sum += u.MiddlePageNumber()
		}
	}

	return sum
}

func parseRules(input string) PrioMap {
	m := PrioMap{}

	for line := range slices.Values(strings.Split(input, "\n")) {
		var left, right int
		_, err := fmt.Sscanf(line, "%d|%d", &left, &right)
		if err != nil {
			panic(fmt.Errorf("error parsing line \"%s\": %w", line, err))
		}

		if _, present := m[left]; !present {
			m[left] = []int{}
		}

		if _, present := m[right]; !present {
			m[right] = []int{}
		}

		m[left] = append(m[left], right)
	}

	return m
}

func parseUpdates(input string) []Update {
	var updates []Update
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		var u Update
		digits := strings.Split(line, ",")
		for _, d := range digits {
			num, err := strconv.Atoi(d)
			if err != nil {
				panic(fmt.Errorf("error parsing %s: %w", d, err))
			}
			u = append(u, num)
		}

		updates = append(updates, u)
	}

	return updates
}

type numeric interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func abs[T numeric](a T) T {
	if a < 0 {
		return a * -1
	}

	return a
}

func floor[T numeric](a T) int {
	return int(math.Floor(float64(a)))
}

func lastIndex[T any](slice []T) int {
	return len(slice) - 1
}
