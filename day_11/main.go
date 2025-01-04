package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"time"
)

type Stones struct {
	m map[int]int
}

func (s *Stones) Blink() {
	newCache := map[int]int{}

	for k, v := range s.m {
		switch {
		case k == 0:
			newCache[1] += v

		case digits(k)%2 == 0:
			left, right := split(k)
			newCache[left] += v
			newCache[right] += v

		default:
			newCache[k*2024] += v
		}
	}

	s.m = newCache
}

func (s *Stones) Count() int {
	var total int
	for _, v := range s.m {
		total += v
	}
	return total
}

func main() {
	stones := parseInput(os.Stdin)

	start := time.Now()
	fmt.Println("answer part one =", partOne(stones))
	fmt.Printf("part one took %+v\n", time.Since(start))

	start = time.Now()
	fmt.Println("answer part two =", partTwo(stones))
	fmt.Printf("part two took %+v", time.Since(start))
}

func blink(num int) []int {
	switch {
	case num == 0:
		return []int{1}

	case digits(num)%2 == 0:
		left, right := split(num)
		return []int{left, right}

	default:
		return []int{num * 2024}
	}
}

func digits(num int) int {
	return int(math.Floor(math.Log10(float64(num)))) + 1
}
func split(num int) (int, int) {
	half := digits(num) / 2
	divisor := int(math.Pow(10, float64(half)))
	return num / divisor, num % divisor
}

func partOne(s Stones) int {
	for i := 0; i < 25; i++ {
		s.Blink()
	}

	return s.Count()
}

func partTwo(s Stones) int {
	// blinking 50 more times since we already blinked 25 times for part one
	for i := 0; i < 50; i++ {
		s.Blink()

	}

	return s.Count()
}

func parseInput(input io.Reader) Stones {
	var (
		scanner = bufio.NewScanner(input)
		res     = Stones{
			m: make(map[int]int),
		}
	)

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		val := scanner.Text()
		num, err := strconv.Atoi(val)
		if err != nil {
			panic(fmt.Errorf("error converting %s to int: %w", val, err))
		}

		res.m[num] = 1
	}

	return res
}
