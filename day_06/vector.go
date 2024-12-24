package main

type Vector struct {
	X int
	Y int
}

func (v Vector) Add(b Vector) Vector {
	return Vector{
		X: v.X + b.X,
		Y: v.Y + b.Y,
	}
}
