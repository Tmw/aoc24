package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

var (
	DirectionUp    = Vector{X: 0, Y: -1}
	DirectionRight = Vector{X: +1, Y: 0}
	DirectionDown  = Vector{X: 0, Y: +1}
	DirectionLeft  = Vector{X: -1, Y: 0}
)

func vectorForInstruction(instr Instruction) Vector {
	switch instr {
	case InstructionUp:
		return DirectionUp
	case InstructionRight:
		return DirectionRight
	case InstructionDown:
		return DirectionDown
	case InstructionLeft:
		return DirectionLeft
	}

	panic(fmt.Errorf("invalid instruction: %s", string(instr)))
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

type Robot struct {
	Vector
}

func (r Robot) clone() Robot {
	return Robot{r.Vector}
}

func (r *Robot) Move(to Vector) {
	r.Vector = to
}

type Tile rune

const (
	TileOpen       = '.'
	TileNonMovable = '#'
	TileMovable    = 'O'
	TileRobot      = '@'

	// for pt 2
	TileMovableLeftHalf  = '['
	TileMovableRightHalf = ']'
)

type Arena struct {
	width  int
	height int

	tiles []Tile
}

func (a *Arena) clone() Arena {
	return Arena{
		width:  a.width,
		height: a.height,
		tiles:  slices.Clone(a.tiles),
	}
}

func (a *Arena) Print(robot Robot) {
	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			loc := Vector{X: x, Y: y}
			tile := a.TileAt(loc)

			if robot.Vector == loc {
				tile = TileRobot
			}

			fmt.Printf("%s", string(tile))
		}
		fmt.Printf("\n")
	}
}

func (a *Arena) FindAllTilesOfType(typ Tile) []Vector {
	var res []Vector
	for idx, tile := range a.tiles {
		if tile == typ {
			x := idx % a.width
			y := idx / a.width
			res = append(res, Vector{X: x, Y: y})
		}
	}

	return res
}

func (a *Arena) TileAt(pos Vector) Tile {
	return a.tiles[pos.Y*a.width+pos.X]
}

func (a *Arena) SetTileAt(pos Vector, t Tile) {
	a.tiles[pos.Y*a.width+pos.X] = t
}

func (a *Arena) InBounds(pos Vector) bool {
	return pos.X > 0 && pos.X < a.width &&
		pos.Y >= 0 && pos.Y <= a.height
}

func (a *Arena) IsTileMovableInDirection(pos Vector, direction Vector) bool {
	var (
		newPos = pos.Add(direction)
		tile   = a.TileAt(newPos)
	)

	if tile == TileOpen {
		return true
	}

	if tile == TileNonMovable {
		return false
	}

	if tile == TileMovable {
		return a.IsTileMovableInDirection(newPos, direction)
	}

	var (
		directionLeftOrRight = direction == DirectionLeft || direction == DirectionRight
		directionUpOrDown    = direction == DirectionUp || direction == DirectionDown
		leftOrRightBoxHalf   = tile == TileMovableLeftHalf || tile == TileMovableRightHalf
	)

	if directionLeftOrRight && leftOrRightBoxHalf {
		// with double width boxes, pushing left or right isn't different than
		// pushing the single boxes left and right.
		return a.IsTileMovableInDirection(newPos, direction)
	}

	if directionUpOrDown && leftOrRightBoxHalf {
		// ok this might be more weird?
		//
		// if this is the left part find the right part (to the right of this one..)
		// return a.IsTileMovableInDirection(left part newpos) && a.IsTileMovableInDirection(right part newpos)..
	}

	panic("reached invalid state")
}

func (a *Arena) MoveTileInDirection(pos Vector, direction Vector) {
	var (
		newPos = pos.Add(direction)
		tile   = a.TileAt(newPos)
	)

	if tile == TileNonMovable {
		return
	}

	if tile == TileOpen {
		a.SetTileAt(newPos, TileMovable)
		a.SetTileAt(pos, TileOpen)
	}

	if tile == TileMovable {
		a.MoveTileInDirection(newPos, direction)
		// re-check the tile at the newPos
		// if now movable, move to it.
		//
		// this ensures we can move entire chains of connected boxes.
		tile = a.TileAt(newPos)
		if tile == TileOpen {
			a.SetTileAt(newPos, TileMovable)
			a.SetTileAt(pos, TileOpen)
		}
	}
}

type Instruction rune

const (
	InstructionUp    = Instruction('^')
	InstructionRight = Instruction('>')
	InstructionDown  = Instruction('v')
	InstructionLeft  = Instruction('<')
)

func main() {
	arena, robot, instructions := parseInput(os.Stdin)
	arenaB, robotB := arena.clone(), robot.clone()

	start := time.Now()
	fmt.Println("answer part one =", partOne(arena, robot, instructions))
	fmt.Printf("part one took %+v\n", time.Since(start))

	applyWidening(&arenaB, &robotB)
	arenaB.Print(robotB)

	start = time.Now()
	fmt.Println("answer part two =", partTwo(arenaB, robotB, instructions))
	fmt.Printf("part two took %+v", time.Since(start))
}

func applyWidening(a *Arena, r *Robot) {
	const wideningFactor = 2
	newTiles := make([]Tile, 0, a.width*wideningFactor*a.height)

	for _, t := range a.tiles {
		if t == TileOpen || t == TileNonMovable {
			newTiles = append(newTiles, []Tile{t, t}...)
			continue
		}

		if t == TileMovable {
			newTiles = append(newTiles, []Tile{TileMovableLeftHalf, TileMovableRightHalf}...)
			continue
		}
	}

	a.width *= wideningFactor
	a.tiles = newTiles
	r.Vector.X *= wideningFactor
}

func runInstructions(arena Arena, robot Robot, instructions []Instruction) {
	for _, instr := range instructions {
		direction := vectorForInstruction(instr)

		newPos := robot.Vector.Add(direction)
		if !arena.InBounds(newPos) {
			continue
		}

		tile := arena.TileAt(newPos)
		if tile == TileNonMovable {
			continue
		}

		if tile == TileOpen {
			robot.Move(newPos)
			continue
		}

		if tile == TileMovable {
			if !arena.IsTileMovableInDirection(newPos, direction) {
				continue
			}

			arena.MoveTileInDirection(newPos, direction)
			robot.Move(newPos)
		}
	}
}

func partOne(arena Arena, robot Robot, instructions []Instruction) int {
	runInstructions(arena, robot, instructions)

	boxPositions := arena.FindAllTilesOfType(TileMovable)
	sum := 0
	for _, p := range boxPositions {
		sum += p.Y*100 + p.X
	}

	return sum
}

func partTwo(arena Arena, robot Robot, instructions []Instruction) int {
	runInstructions(arena, robot, instructions)

	// from the example it looked like only the distance from the left half was measured,
	// ignoring the "closest edge" wording.
	boxPositions := arena.FindAllTilesOfType(TileMovableLeftHalf)
	sum := 0
	for _, p := range boxPositions {
		sum += p.Y*100 + p.X
	}

	return sum
}

type ParseType int

const (
	ParseTypeMap = iota
	ParseTypeInstructions
)

func parseInput(input io.Reader) (Arena, Robot, []Instruction) {
	var (
		scanner      = bufio.NewScanner(input)
		arena        Arena
		robot        Robot
		instructions []Instruction
		parseType    = ParseTypeMap
	)

	for scanner.Scan() {
		line := scanner.Text()

		// found the double linebreak,
		// now parsing instructions
		if parseType == ParseTypeMap && line == "" {
			parseType = ParseTypeInstructions
			continue
		}

		// some janky parsing code, oops.
		// but this parses the map tiles, sets width and height and
		// finds the initial position for the robot.
		if parseType == ParseTypeMap {
			arena.width = 0
			for _, c := range line {
				tile := Tile(c)

				// detect robot, parse location but store tile as open
				if tile == TileRobot {
					robot.X = arena.width
					robot.Y = arena.height
					tile = TileOpen
				}

				arena.tiles = append(arena.tiles, tile)
				arena.width++
			}

			arena.height++
			continue
		}

		// parsing of the instructions is straight forward
		if parseType == ParseTypeInstructions {
			for _, c := range line {
				instructions = append(instructions, Instruction(c))
			}
		}
	}

	return arena, robot, instructions
}
