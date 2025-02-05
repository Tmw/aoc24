package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
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
	objects := make([]*Object, 0, len(w.objects))
	for _, o := range w.objects {
		objects = append(objects, &Object{
			width:  o.width,
			height: o.height,
			pos:    o.pos,
			typ:    o.typ,
		})
	}

	return World{
		width:   w.width,
		height:  w.height,
		objects: objects,
	}
}

func (w *World) FindObjectsOfType(typ ObjTyp) []*Object {
	var res []*Object
	for _, obj := range w.objects {
		if obj.typ == typ {
			res = append(res, obj)
		}
	}

	return res
}

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

	world2, robot2 := world.Clone(), robot.Clone()

	start := time.Now()
	fmt.Println("answer part one =", getAnswer(&world, &robot, instructions))
	fmt.Printf("part one took %+v\n", time.Since(start))

	applyWidening(&world2, &robot2)

	start = time.Now()
	fmt.Println("answer part two =", getAnswer(&world2, &robot2, instructions))
	fmt.Printf("part two took %+v", time.Since(start))
}

func applyWidening(w *World, r *Robot) {
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
		direction := vectorForInstruction(instr)
		newPos := robot.Vector.Add(direction)

		if !world.InBounds(newPos) {
			continue
		}

		newRobot := &Object{pos: newPos, width: 1, height: 1}
		if !world.CanMove(newRobot, direction) {
			continue
		}

		objsInPath := world.FindCollisions(newRobot, direction, []*Object{})
		for _, o := range unique(objsInPath) {
			o.MoveInDirection(direction)
		}

		robot.Move(newPos)
	}
}

func getAnswer(world *World, robot *Robot, instructions []Instruction) int {
	runInstructions(world, robot, instructions)
	boxPositions := world.FindObjectsOfType(ObjTypeBox)
	sum := 0
	for _, o := range boxPositions {
		sum += o.pos.Y*100 + o.pos.X
	}

	return sum
}

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

func unique[T comparable](input []T) []T {
	unique := make(map[T]struct{})
	for _, item := range input {
		unique[item] = struct{}{}
	}

	return slices.Collect(maps.Keys(unique))
}
