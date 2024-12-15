package main

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func assert(t *testing.T, a, b any, msg string) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf(fmt.Sprintf("assertion failed: %s. Expected %+v, got %+v", msg, b, a))
	}
}

func TestParseGrid(t *testing.T) {
	g := ParseGrid("abc\ndef\nghi\n")
	assert(t, g.width, 3, "grid width incorrect")
	assert(t, g.height, 3, "grid height incorrect")
	assert(t, g.CharAt(0, 0), byte('a'), "incorrect char")
	assert(t, g.CharAt(1, 0), byte('b'), "incorrect char")
	assert(t, g.CharAt(0, 1), byte('d'), "incorrect char")
	assert(t, g.CharAt(2, 1), byte('f'), "incorrect char")
	assert(t, g.CharAt(2, 2), byte('i'), "incorrect char")
}

func TestHorizontals(t *testing.T) {
	g := ParseGrid("abc\ndef\nghi\n")
	h := slices.Collect(g.Horizontals())
	assert(t, len(h), 3, "incorrect length")
	assert(t, h[0], []byte("abc"), "incorrect first line")
	assert(t, h[1], []byte("def"), "incorrect second line")
	assert(t, h[2], []byte("ghi"), "incorrect third line")
}

func TestVerticals(t *testing.T) {
	g := ParseGrid("abc\ndef\nghi\n")
	h := slices.Collect(g.Verticals())
	assert(t, len(h), 3, "incorrect length")
	assert(t, h[0], []byte("adg"), "incorrect first line")
	assert(t, h[1], []byte("beh"), "incorrect second line")
	assert(t, h[2], []byte("cfi"), "incorrect third line")
}

func TestDiagonal(t *testing.T) {
	g := ParseGrid("abc\ndef\nghi\n")

	t.Run("left-to-right", func(t *testing.T) {
		assert(t, g.Diagonal(0, 0, +1), []byte("aei"), "wrong diagonal")
		assert(t, g.Diagonal(1, 0, +1), []byte("bf"), "wrong diagonal")
		assert(t, g.Diagonal(2, 0, +1), []byte("c"), "wrong diagonal")
		assert(t, g.Diagonal(0, 1, +1), []byte("dh"), "wrong diagonal")
		assert(t, g.Diagonal(0, 2, +1), []byte("g"), "wrong diagonal")
	})

	t.Run("right-to-left", func(t *testing.T) {
		fmt.Printf("diagonal = %s\n", g.Diagonal(2, 2, -1))
		assert(t, g.Diagonal(2, 0, -1), []byte("ceg"), "wrong diagonal")
		assert(t, g.Diagonal(1, 0, -1), []byte("bd"), "wrong diagonal")
		assert(t, g.Diagonal(0, 0, -1), []byte("a"), "wrong diagonal")
		assert(t, g.Diagonal(2, 1, -1), []byte("fh"), "wrong diagonal")
		assert(t, g.Diagonal(2, 2, -1), []byte("i"), "wrong diagonal")
	})
}

func TestDiagonals(t *testing.T) {
	g := ParseGrid("abc\ndef\nghi\n")
	h := slices.Collect(g.Diagonals())

	// a b c
	// d e f
	// g h i
	//
	// bcomes =
	//  left to right       right to left
	// - g                   i
	// - d h                 f h
	// - a e i               c e g
	// - b f                 b d
	// - c                   a
	assert(t, len(h), 10, "incorrect length")
	assert(t, h[0], []byte("g"), "incorrect first line")
	assert(t, h[1], []byte("dh"), "incorrect second line")
	assert(t, h[2], []byte("aei"), "incorrect third line")
	assert(t, h[3], []byte("bf"), "incorrect third line")
	assert(t, h[4], []byte("c"), "incorrect third line")
	assert(t, h[5], []byte("i"), "incorrect first line")
	assert(t, h[6], []byte("fh"), "incorrect second line")
	assert(t, h[7], []byte("ceg"), "incorrect third line")
	assert(t, h[8], []byte("bd"), "incorrect third line")
	assert(t, h[9], []byte("a"), "incorrect third line")
}

func TestSubgrid(t *testing.T) {
	// entire grid:
	//
	// a b c d e
	// f g h i j
	// k l m n o
	// p q r s t
	// u v w x y

	g := ParseGrid("abcde\nfghij\nklmno\npqrst\nuvwxy\n")

	// subgrid (0,0,3,3):
	// [a b c] d e
	// [f g h] i j
	// [k l m] n o
	//  p q r  s t
	//  u v w  x y
	sg := g.Subgrid(0, 0, 3, 3)
	assert(t, sg.width, 3, "subgrid width incorrect")
	assert(t, sg.height, 3, "subgrid height incorrect")
	assert(t, sg.CharAt(0, 0), byte('a'), "incorrect char in subgrid")
	assert(t, sg.CharAt(1, 1), byte('g'), "incorrect char in subgrid")
	assert(t, sg.CharAt(2, 2), byte('m'), "incorrect char in subgrid")

	// subgrid (2,2,3,3):
	// a b  c d e
	// f g  h i j
	// k l [m n o]
	// p q [r s t]
	// u v [w x y]
	sg = g.Subgrid(2, 2, 3, 3)
	assert(t, sg.width, 3, "subgrid width incorrect")
	assert(t, sg.height, 3, "subgrid height incorrect")
	assert(t, sg.CharAt(0, 0), byte('m'), "incorrect char in subgrid")
	assert(t, sg.CharAt(1, 1), byte('s'), "incorrect char in subgrid")
	assert(t, sg.CharAt(2, 2), byte('y'), "incorrect char in subgrid")

	// subgrid diagonals
	assert(t, sg.Diagonal(0, 0, 1), []byte("msy"), "incorrect subgrid diagonal")
	assert(t, sg.Diagonal(2, 0, -1), []byte("osw"), "incorrect subgrid diagonal")
}

func TestCombineIter(t *testing.T) {
	iter1 := slices.Values([]int{1, 2, 3})
	iter2 := slices.Values([]int{4, 5, 6})

	combined := slices.Collect(CombineIter(iter1, iter2))
	assert(t, combined, []int{1, 2, 3, 4, 5, 6}, "combine not combining")
}
