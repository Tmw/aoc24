package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Op int

const (
	OpAdd Op = iota
	OpMul
	OpCat // =^.^=
)

var AvailableOps = []Op{
	OpAdd,
	OpMul,
	OpCat,
}

type Equation struct {
	Sum   int
	Parts []int
}

func main() {
	equations := parseInput(os.Stdin)

	p1Stop := profile("part one")
	fmt.Println("answer part one =", solve(equations, AvailableOps[0:2]))
	p1Stop()

	p2Stop := profile("part two")
	fmt.Println("answer part two =", solve(equations, AvailableOps))
	p2Stop()
}

func solve(equations []Equation, availableOps []Op) int {
	sum := 0
	for _, eq := range equations {
		if check(eq.Sum, 0, eq.Parts, availableOps) {
			sum += eq.Sum
		}
	}

	return sum
}

// intCat concatenates two integers by multiplying A by 10 N times
// where N is the number of digits in B.
//
// Example:
// 10 cat 9 is done in the following way:
// 10 * (10 * 1) (*1 as 9 is 1 digit) = 100
// 100 + 9 = 109
func intCat(a, b int) int {
	for bCopy := b; bCopy > 0; bCopy /= 10 {
		a *= 10
	}

	return a + b
}

func check(target, total int, nums []int, ops []Op) bool {
	if len(nums) == 0 {
		return total == target
	}

	if total > target {
		return false
	}

	var (
		num  = nums[0]
		rest = nums[1:]
	)

	for _, op := range ops {
		switch op {
		case OpMul:
			if check(target, total*num, rest, ops) {
				return true
			}

		case OpAdd:
			if check(target, total+num, rest, ops) {
				return true
			}

		case OpCat:
			if check(target, intCat(total, num), rest, ops) {
				return true
			}
		}
	}

	return false
}

func parseInput(input io.Reader) []Equation {
	var (
		output  []Equation
		scanner = bufio.NewScanner(input)
	)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		sum, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Errorf("error parsing integer %+v: %w", parts[0], err))
		}

		var l Equation
		l.Sum = sum

		for _, part := range strings.Split(parts[1], " ") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			p, err := strconv.Atoi(part)
			if err != nil {
				panic(fmt.Errorf("error parsing integer '%+v': %w", part, err))
			}

			l.Parts = append(l.Parts, p)
		}

		output = append(output, l)
	}

	return output
}

func profile(label string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("profile %s took %+v\n", label, time.Since(start))
	}
}
