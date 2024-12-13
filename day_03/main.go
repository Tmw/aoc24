package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type InstructionType string

const (
	InstructionTypeMul  InstructionType = "mul"
	InstructionTypeDo   InstructionType = "do"
	InstructionTypeDont InstructionType = "dont"
)

type Instruction struct {
	Position int
	Typ      InstructionType
	Expr     MultiplicationExpression
}

func (i Instruction) String() string {
	switch i.Typ {
	case InstructionTypeMul:
		return fmt.Sprintf("mul(%d,%d)", i.Expr.A, i.Expr.B)

	case InstructionTypeDo:
		return "do()"

	case InstructionTypeDont:
		return "don't()"

	default:
		return "<unknown>"
	}
}

type MultiplicationExpression struct {
	A int
	B int
}

func (m MultiplicationExpression) Mul() int {
	return m.A * m.B
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("unable to read input: %v", err)
	}

	instructions := parseWithoutRegex(string(input))

	// calculate part one
	sum := 0
	for ins := range slices.Values(instructions) {
		if ins.Typ == InstructionTypeMul {
			sum += ins.Expr.Mul()
		}
	}

	fmt.Println("part one = ", sum) // should be 170778545

	// calculate part two
	sum = 0
	enabled := true
	for ins := range slices.Values(instructions) {
		switch ins.Typ {
		case InstructionTypeDo:
			enabled = true
		case InstructionTypeDont:
			enabled = false
		case InstructionTypeMul:
			if !enabled {
				continue
			}

			sum += ins.Expr.Mul()
		}
	}

	fmt.Println("part two = ", sum) // 82868252
}

func parseWithoutRegex(input string) []Instruction {
	var (
		offset = 0

		// buffer size equals max length of instruction: mul(ddd,ddd)
		bufferSize   = 12
		instructions []Instruction
	)

	for offset < len(input) {
		chunk := input[offset:min(offset+bufferSize, len(input))]

		// match literal do() instruction
		if strings.HasPrefix(chunk, "do()") {
			ins := Instruction{
				Typ:      InstructionTypeDo,
				Position: offset,
			}
			instructions = append(instructions, ins)
			offset += 4
			continue
		}

		// match literal don't() instruction
		if strings.HasPrefix(chunk, "don't()") {
			ins := Instruction{
				Typ:      InstructionTypeDont,
				Position: offset,
			}
			instructions = append(instructions, ins)
			offset += 4
			continue
		}

		// match mul(ddd,ddd) instruction
		closingParenIdx := strings.Index(chunk, ")")
		if strings.HasPrefix(chunk, "mul(") && closingParenIdx > -1 {
			token := chunk[:closingParenIdx+1]
			var expr MultiplicationExpression
			_, err := fmt.Sscanf(token, "mul(%d,%d)", &expr.A, &expr.B)
			if err != nil {
				log.Printf("error trying to scan %s at offset %d: %v", token, offset, err)
				offset++
				continue
			}

			ins := Instruction{
				Typ:      InstructionTypeMul,
				Position: offset,
				Expr:     expr,
			}
			instructions = append(instructions, ins)
			offset += len(token)
			continue
		}

		offset++
	}

	return instructions
}
