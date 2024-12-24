package main

import (
	"fmt"
)

type Heading int

const (
	HeadingNorth Heading = iota
	HeadingEast
	HeadingSouth
	HeadingWest
)

func (h Heading) Vector() Vector {
	switch h {
	case HeadingNorth:
		return Vector{X: 0, Y: -1}

	case HeadingEast:
		return Vector{X: 1, Y: 0}

	case HeadingSouth:
		return Vector{X: 0, Y: 1}

	case HeadingWest:
		return Vector{X: -1, Y: 0}
	}

	panic(fmt.Errorf("invalid heading: %d", h))
}

func HeadingFromVector(v Vector) Heading {
	switch {
	case v.X == 0 && v.Y < 0:
		return HeadingNorth

	case v.X > 0 && v.Y == 0:
		return HeadingEast

	case v.X == 0 && v.Y > 0:
		return HeadingSouth

	case v.X < 0 && v.Y == 0:
		return HeadingWest
	}

	panic(fmt.Errorf("unable to determine heading from vector: %+v", v))
}

func (h Heading) RotateClockwise() Heading {
	if h == HeadingWest {
		return HeadingNorth
	}

	return h + 1
}
