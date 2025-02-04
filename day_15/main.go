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

type ObjTyp int

const (
	ObjTypeWall = ObjTyp(iota)
	ObjTypeBox
)

type Object struct {
	width  int
	height int
	pos    Vector
	typ    ObjTyp
}

func (o *Object) MoveInDirection(direction Vector) {
	o.pos = o.pos.Add(direction)
}

type World struct {
	width   int
	height  int
	objects []*Object
}

func (w World) Print(r *Robot) {
	// make empty canvas
	res := make([]rune, 0, w.width*w.height)
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			res = append(res, '.')
		}
	}

	// draw in the objects
	for _, obj := range w.objects {
		var (
			start = obj.pos.Y*w.width + obj.pos.X
			end   = min(obj.pos.Y*w.width+obj.pos.X+obj.width, len(res))
		)

		for idx := start; idx < end; idx++ {
			switch obj.typ {
			case ObjTypeWall:
				res[idx] = '#'
			case ObjTypeBox:
				if obj.width == 1 {
					res[idx] = 'O'
					continue
				}

				if idx == end-1 {
					res[idx] = ']'
					continue
				}

				if idx == start {
					res[idx] = '['
					continue
				}
			}
		}

	}

	// draw in the robot
	res[r.Y*w.width+r.X] = '@'

	// draw each line
	for i, c := range res {
		if i > 0 && i%w.width == 0 {
			fmt.Print("\n")
		}

		fmt.Print(string(c))
	}

	fmt.Println()
}

func (w World) InBounds(pos Vector) bool {
	return pos.X >= 0 && pos.Y >= 0 && pos.X <= w.width && pos.Y <= w.height
}

func (w World) FindCollisions(obj1 *Object, direction Vector, ignoring []*Object) []*Object {
	var collisions []*Object
	for _, obj2 := range w.objects {
		// no comparison with self
		if obj1 == obj2 {
			continue
		}

		// skip comparing against objects in ignoring list
		if slices.Contains(ignoring, obj2) {
			continue
		}

		// do objects intersect?
		intersect := obj1.pos.X < obj2.pos.X+obj2.width &&
			obj1.pos.X+obj1.width > obj2.pos.X &&
			obj1.pos.Y < obj2.pos.Y+obj2.height &&
			obj1.pos.Y+obj1.height > obj2.pos.Y

		if intersect {
			collisions = append(collisions, obj2)
		}
	}

	for _, colliding := range collisions {
		newPos := colliding.pos.Add(direction)
		probe := &Object{width: colliding.width, height: colliding.height, pos: newPos}
		moreCollissions := w.FindCollisions(probe, direction, []*Object{colliding})
		collisions = append(collisions, moreCollissions...)
	}

	return collisions
}

// CanMove returns true if in the obj van move in direction
// either because it find an empty space or the objects it collides with
// are in turn movable too.
func (w World) CanMove(obj *Object, direction Vector) bool {
	collisions := w.FindCollisions(obj, direction, []*Object{})

	for _, col := range collisions {
		if col.typ == ObjTypeWall {
			return false
		}
	}

	return true
}

func (w *World) Clone() World {
	return World{
		width:   w.width,
		height:  w.height,
		objects: slices.Clone(w.objects),
	}
}

// func (a *Arena) FindAllTilesOfType(typ Tile) []Vector {
// 	var res []Vector
// 	for idx, tile := range a.tiles {
// 		if tile == typ {
// 			x := idx % a.width
// 			y := idx / a.width
// 			res = append(res, Vector{X: x, Y: y})
// 		}
// 	}
//
// 	return res
// }

type Robot struct {
	Vector
}

func (r Robot) Clone() Robot {
	return Robot{r.Vector}
}

func (r *Robot) Move(to Vector) {
	r.Vector = to
}

const (
	TileOpen       = '.'
	TileNonMovable = '#'
	TileMovable    = 'O'
	TileRobot      = '@'
)

type Instruction rune

const (
	InstructionUp    = Instruction('^')
	InstructionRight = Instruction('>')
	InstructionDown  = Instruction('v')
	InstructionLeft  = Instruction('<')
)

func main() {
	world, robot, instructions := parseInput(os.Stdin)
	_ = instructions

	world.Print(&robot)
	// world2, robot2 := world.Clone(), robot.Clone()

	start := time.Now()
	fmt.Println("answer part one =", partOne(&world, &robot, instructions))
	fmt.Printf("part one took %+v\n", time.Since(start))

	// applyWidening(&world2, &robot2)
	// world2.Print(robot2)

	// start := time.Now()
	// fmt.Println("answer part two =", partTwo(arenaB, robotB, instructions))
	// fmt.Printf("part two took %+v", time.Since(start))
}

func applyWidening(w *World, r *Robot) {
	_ = r
	const wideningFactor = 2

	w.width *= wideningFactor
	for idx := range w.objects {
		w.objects[idx].pos.X *= wideningFactor
		w.objects[idx].width *= wideningFactor
	}

	r.Vector.X *= wideningFactor
}

func runInstructions(world *World, robot *Robot, instructions []Instruction) {
	for _, instr := range instructions {
		fmt.Println("-------------------------------")
		world.Print(robot)
		fmt.Println("instruction = ", string(instr))

		// time.Sleep(200 * time.Millisecond)

		direction := vectorForInstruction(instr)
		newPos := robot.Vector.Add(direction)

		if !world.InBounds(newPos) {
			continue
		}

		// TODO: Add method for moving the robot? Expose through Vector?
		// should Robot be just another Object in the world? :idea:
		newRobot := &Object{pos: newPos, width: 1, height: 1}
		if !world.CanMove(newRobot, direction) {
			fmt.Println("unable to move, continuing..")
			continue
		}

		// TODO: Let's make it work, then improve:
		// - perhaps CanMove can return a bool and a list of objects to move.
		// - then calling move individually on them does the trick?
		objsInPath := world.FindCollisions(newRobot, direction, []*Object{})
		for _, o := range objsInPath {
			o.MoveInDirection(direction)
		}

		robot.Move(newPos)
	}
}

func partOne(world *World, robot *Robot, instructions []Instruction) int {
	applyWidening(world, robot)
	runInstructions(world, robot, instructions)
	fmt.Println("final state of the world:")
	world.Print(robot)

	// boxPositions := world.FindAllTilesOfType(TileMovable)
	// sum := 0
	// for _, p := range boxPositions {
	// 	sum += p.Y*100 + p.X
	// }

	var sum int
	return sum
}

//
// func partTwo(arena Arena, robot Robot, instructions []Instruction) int {
// 	runInstructions(arena, robot, instructions)
//
// 	// from the example it looked like only the distance from the left half was measured,
// 	// ignoring the "closest edge" wording.
// 	boxPositions := arena.FindAllTilesOfType(TileMovableLeftHalf)
// 	sum := 0
// 	for _, p := range boxPositions {
// 		sum += p.Y*100 + p.X
// 	}
//
// 	return sum
// }

type ParseType int

const (
	ParseTypeMap = iota
	ParseTypeInstructions
)

func parseInput(input io.Reader) (World, Robot, []Instruction) {
	var (
		scanner      = bufio.NewScanner(input)
		world        World
		robot        Robot
		instructions []Instruction
		parseType    = ParseTypeMap
		lineNo       int
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
			for idx, c := range line {
				switch c {
				case '#':
					world.objects = append(world.objects, &Object{
						width:  1,
						height: 1,
						typ:    ObjTypeWall,
						pos: Vector{
							X: idx,
							Y: lineNo,
						},
					})

				case 'O':
					world.objects = append(world.objects, &Object{
						width:  1,
						height: 1,
						typ:    ObjTypeBox,
						pos: Vector{
							X: idx,
							Y: lineNo,
						},
					})

				case '@':
					robot.X = idx
					robot.Y = lineNo
				}

				world.width = idx + 1
			}

			lineNo++
			world.height = lineNo

			continue
		}

		// parsing of the instructions is straight forward
		if parseType == ParseTypeInstructions {
			for _, c := range line {
				instructions = append(instructions, Instruction(c))
			}
		}
	}

	return world, robot, instructions
}
