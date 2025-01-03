package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {
	s := parseInput(os.Stdin)

	start := time.Now()
	fmt.Println("answer part one =", partOne(s))
	fmt.Printf("part one took %+v\n", time.Since(start))

	start = time.Now()
	fmt.Println("answer part two =", partTwo(s))
	fmt.Printf("part two took %+v", time.Since(start))
}

type Stone struct {
	Val int
}

func (s *Stone) Digits() int {
	return int(math.Floor(math.Log10(float64(s.Val)))) + 1
}

type Stones struct {
	list *list.List
}

func (s *Stones) Blink() {
	for node := s.list.Front(); node != nil; node = node.Next() {
		stone, ok := node.Value.(*Stone)
		if !ok {
			continue
		}

		switch {
		case stone.Val == 0:
			stone.Val = 1

		case stone.Digits()%2 == 0:
			half := stone.Digits() / 2
			divisor := int(math.Pow(10, float64(half)))
			left, right := stone.Val/divisor, stone.Val%divisor

			stone.Val = left

			// Adding a new node after the current node, and forwarding the iterator
			// by one so we don't apply the blink to the new element next.
			node = s.list.InsertAfter(&Stone{Val: right}, node)

		default:
			stone.Val *= 2024
		}
	}
}

func (s *Stones) Print() {
	for node := s.list.Front(); node != nil; node = node.Next() {
		s, ok := node.Value.(*Stone)
		if !ok {
			continue
		}

		fmt.Printf("%d ", s.Val)
	}

	fmt.Println()
}

func partOne(s Stones) int {
	for i := 0; i < 25; i++ {
		s.Blink()
	}

	return s.list.Len()
}

func partTwo(s Stones) int {
	// blinking 50 more times since we already blinked 25 times for part one
	for i := 0; i < 50; i++ {
		s.Blink()

		fmt.Printf("after blinking %d times we have %d stones.\n", i+25, s.list.Len())
	}

	return s.list.Len()
}

func parseInput(input io.Reader) Stones {
	var (
		scanner = bufio.NewScanner(input)
		res     = Stones{list: list.New()}
	)

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		val := scanner.Text()
		num, err := strconv.Atoi(val)
		if err != nil {
			panic(fmt.Errorf("error converting %s to int: %w", val, err))
		}

		res.list.PushBack(&Stone{Val: num})
	}

	return res
}
