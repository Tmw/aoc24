package main

import (
	"math"
	"slices"

	"github.com/tmw/go-prioqueue"
)

type Candidate struct {
	loc Vector

	// total cost of this candidate. This contains
	// distance from the start, heuristic to finish and optional penalty
	cost int

	// distance from start
	dist int

	// via what candidate we got here.
	// will update if we reach the same node through a shorter path.
	// is used to backtrack into a path once we reach the finish.
	via *Candidate

	// what direction did we come from when we landed here
	dir Vector
}

type PathFinder struct {
	PathFinderOpts
	visited map[Vector]Candidate
	queue   prioqueue.PrioQueue[Candidate, Vector]
}

type PathFinderOpts struct {
	NeighboursFn    func(l Vector) []Vector
	HeuristicFn     func(l Vector) int
	ReachedFinishFn func(l Vector) bool
}

func NewPathFinder(opts PathFinderOpts) PathFinder {
	return PathFinder{
		PathFinderOpts: opts,
		visited:        make(map[Vector]Candidate),
		queue: prioqueue.NewPrioQueue(
			compareCandidate,
			hashCandidate,
		),
	}
}
func (p *PathFinder) Path(start Vector) (int, []Vector) {
	initialCost := p.HeuristicFn(start)
	p.queue.Push(Candidate{
		loc:  start,
		cost: initialCost,
		dist: 0,
		via:  nil,
		dir:  DirectionEast,
	})

	for {
		c, more := p.queue.Pop()
		if !more {
			break
		}

		if p.ReachedFinishFn(c.loc) {
			return c.cost, backtrack(c)
		}

		for _, n := range p.NeighboursFn(c.loc) {
			dist := c.dist + 1

			if c.loc.Add(c.dir) != n {
				dist += 1000
			}

			cost := p.HeuristicFn(n) + dist
			candidate := Candidate{
				loc:  n,
				dist: dist,
				cost: cost,
				via:  &c,
				dir:  n.Sub(c.loc),
			}

			ok := p.queue.Update(hashCandidate(candidate), func(c Candidate) Candidate {
				if c.cost <= candidate.cost {
					return c
				}

				c.cost = candidate.cost
				c.dist = candidate.dist
				c.via = candidate.via
				c.dir = candidate.dir

				return c
			})

			// if neighbour already on processed list, consider re-adding
			// but only if its distance from start is lower.
			if p, found := p.visited[n]; found && p.cost <= candidate.cost {
				continue
			}

			if !ok {
				// candidate unknown, pushing new one
				p.queue.Push(candidate)
			}

		}

		p.visited[c.loc] = c
	}

	return math.MaxInt, []Vector{}
}

func backtrack(c Candidate) []Vector {
	var path []Vector

	for {
		path = append(path, c.loc)
		if c.via == nil {
			break
		}

		c = *c.via
	}

	slices.Reverse(path)
	return path
}

func hashCandidate(c Candidate) Vector {
	return c.loc
}

func compareCandidate(a, b Candidate) bool {
	return a.cost < b.cost
}
