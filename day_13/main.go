package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

const PartTwoOffset = 10_000_000_000_000

var PartTwoOffsetVector = Vector{
	X: PartTwoOffset,
	Y: PartTwoOffset,
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

type Machine struct {
	Prize   Vector
	ButtonA Vector
	ButtonB Vector
}

func main() {
	machines := parseInput(os.Stdin)

	start := time.Now()
	fmt.Println("answer part one =", cost(machines))
	fmt.Printf("part one took %+v\n", time.Since(start))

	applyPartTwoOffset(machines)
	start = time.Now()
	fmt.Println("answer part two =", cost(machines))
	fmt.Printf("part two took %+v", time.Since(start))
}

func cost(m []Machine) int {
	total := 0

	for _, machine := range m {
		cost, found := findPrizeCramersRule(machine)
		if found {
			total += cost
		}
	}

	return total
}

func applyPartTwoOffset(machines []Machine) {
	for idx := range machines {
		machines[idx].Prize = machines[idx].Prize.Add(PartTwoOffsetVector)
	}
}

// naive approach used for part one: nested for loops
func findPrize(m Machine) (int, bool) {
	for a := 0; a < 100; a++ {
		for b := 0; b < 100; b++ {
			aLoc := Vector{X: m.ButtonA.X * a, Y: m.ButtonA.Y * a}
			bLoc := Vector{X: m.ButtonB.X * b, Y: m.ButtonB.Y * b}
			clawLoc := aLoc.Add(bLoc)

			if clawLoc.X > m.Prize.X || clawLoc.Y > m.Prize.Y {
				// break the innerloop if we've gone too far
				break
			}

			if m.Prize == clawLoc {
				return a*3 + b*1, true
			}
		}
	}

	return math.MaxInt, false
}

func findPrizeCramersRule(m Machine) (int, bool) {
	//x = a * buttonA.X + b * buttonB.X
	//y = a * buttonA.Y + b * buttonB.Y
	det := float64(m.ButtonA.X)*float64(m.ButtonB.Y) - float64(m.ButtonB.X)*float64(m.ButtonA.Y)
	if det == 0 {
		return -1, false
	}

	// | px ax |
	// | py ay |
	da := float64(m.Prize.X)*float64(m.ButtonB.Y) - float64(m.ButtonB.X)*float64(m.Prize.Y)

	// | bx px |
	// | by py |
	db := float64(m.ButtonA.X)*float64(m.Prize.Y) - float64(m.Prize.X)*float64(m.ButtonA.Y)
	a := da / det
	b := db / det

	if math.Trunc(a) != a || math.Trunc(b) != b {
		return -1, false
	}

	return int(a*3 + b*1), true
}

func parseInput(input io.Reader) []Machine {
	var (
		scanner  = bufio.NewScanner(input)
		machines = []Machine{}
	)

	scanner.Split(scanBlocks)
	for scanner.Scan() {
		block := scanner.Text()

		var m Machine
		lines := strings.Split(block, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}

			var x, y int
			if _, err := fmt.Sscanf(line, "Button A: X+%d, Y+%d", &x, &y); err == nil {
				m.ButtonA = Vector{X: x, Y: y}
			}

			if _, err := fmt.Sscanf(line, "Button B: X+%d, Y+%d", &x, &y); err == nil {
				m.ButtonB = Vector{X: x, Y: y}
			}

			if _, err := fmt.Sscanf(line, "Prize: X=%d, Y=%d", &x, &y); err == nil {
				m.Prize = Vector{X: x, Y: y}
			}
		}
		machines = append(machines, m)
	}

	return machines
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func scanBlocks(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), "\n\n"); i >= 0 {
		return i + 1, dropCR(data[0:i]), nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
