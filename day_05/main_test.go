package main

import (
	"testing"
)

func assert(t *testing.T, statement bool, message string) {
	if !statement {
		t.Errorf("assertion failed: %s", message)
	}
}

var priomap = `47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13`

func TestUpdate_MiddlePageNumber(t *testing.T) {
	u := Update{1, 2, 3, 4, 5}
	assert(t, u.MiddlePageNumber() == 3, "incorrect middle page number")
}

func TestPrioMap_Scope(t *testing.T) {
	p := parseRules(priomap)
	updates := []Update{
		[]int{75, 47, 61, 53, 29},
		[]int{97, 61, 53, 29, 13},
		[]int{75, 29, 13},
		[]int{75, 97, 47, 61, 53},
		[]int{61, 13, 29},
		[]int{97, 13, 75, 29, 47},
	}

	for _, u := range updates {
		s := p.Scope(u)
		assert(t, len(s) == len(u), "incorrect scope")
	}
}
