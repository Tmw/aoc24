package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
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

func (e *Equation) Validate(opSeq []Op) bool {
	sum := 0
	for idx, part := range e.Parts {
		if idx == 0 {
			// no need for special operator for first part.
			// simply using addition so the number doesn't get lost.
			sum += part
			continue
		}

		switch opSeq[idx-1] {
		case OpMul:
			sum *= part

		case OpAdd:
			sum += part

		case OpCat:
			combined := fmt.Sprintf("%d%d", sum, part)
			res, err := strconv.Atoi(combined)
			if err != nil {
				panic(fmt.Errorf("error parsing '%s' as int: %w", combined, err))
			}
			sum = res
		}
	}

	return sum == e.Sum
}

func main() {
	equations := parseInput(os.Stdin)
	fmt.Println("answer part one =", solve(equations, AvailableOps[0:2]))
	fmt.Println("answer part two =", solve(equations, AvailableOps))
}

func solve(equations []Equation, availableOps []Op) int {
	var (
		opsPermCache = map[int][][]Op{}
		sum          = 0
	)

	for _, eq := range equations {
		var ops [][]Op
		ops, found := opsPermCache[len(eq.Parts)-1]
		if !found {
			ops = permutations(len(eq.Parts)-1, availableOps)
			opsPermCache[len(eq.Parts)-1] = ops
		}

		for _, opSeq := range ops {
			if eq.Validate(opSeq) {
				sum += eq.Sum
				break
			}
		}
	}

	return sum
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

func permutations(length int, options []Op) [][]Op {
	var (
		n       = len(options)
		total   = pow(len(options), length)
		results [][]Op
	)

	for i := 0; i < total; i++ {
		var permutation []Op
		num := i // Current number in base-n
		for j := 0; j < length; j++ {
			optionIndex := num % n
			permutation = append([]Op{options[optionIndex]}, permutation...)
			num /= n
		}
		results = append(results, permutation)
	}

	return results
}

func pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}
