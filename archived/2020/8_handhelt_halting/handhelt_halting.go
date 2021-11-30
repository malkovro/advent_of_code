package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func readFile(inputFile string) string {
	fmt.Printf("Reading map from: %v \n", inputFile)

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("File reading error", err)
		return ""
	}

	return string(data)
}

type Instruction struct {
	Name  string
	Value int
}

func compileInstruction(str string) Instruction {
	splitted := strings.Split(str, " ")

	i, _ := strconv.Atoi(splitted[1])
	return Instruction{splitted[0], i}
}

func exexuteInstruction(instruction Instruction, index int, acc int) (int, int) {
	switch instruction.Name {
	case "jmp":
		return (index + instruction.Value), acc
	case "acc":
		acc += instruction.Value
	}
	return (index + 1), acc
}

func runProgram(instructions []Instruction) (accumulator int, err string) {
	index := 0
	visitedInstructions := make(map[int]bool)
	_, visited := visitedInstructions[index]

	for !visited {
		visitedInstructions[index] = true
		if index == len(instructions) {
			return
		}
		if index > len(instructions) || index < 0 {
			err = "Jumping to the limbo"
			return
		}
		nextInstruction := instructions[index]
		index, accumulator = exexuteInstruction(nextInstruction, index, accumulator)
		_, visited = visitedInstructions[index]
	}

	err = "Looped on myself"
	return
}

func main() {
	inputFile := os.Args[1]
	code := readFile(inputFile)

	instructions := []Instruction{}

	for _, strInstruction := range strings.Split(code, "\n") {
		instructions = append(instructions, compileInstruction(strInstruction))
	}

	instructionsFixedPotential := make([]Instruction, len(instructions))

	faultyInstructionIndex := -1

	acc, err := runProgram(instructions)
	fmt.Println("Running the program as is: ", err, ". Accumulator holds: ", acc)

	for err != "" {
		switchedSomething := false
		faultyInstructionIndex++
		if faultyInstructionIndex == len(instructions) {
			fmt.Println("Could not save private Ryan! 😭")
			break
		}

		copy(instructionsFixedPotential, instructions)

		if instructions[faultyInstructionIndex].Name == "jmp" {
			switchedSomething = true
			fmt.Println("Tried switching `jmp` to `nop` statement line", faultyInstructionIndex)
			instructionsFixedPotential[faultyInstructionIndex].Name = "nop"
		}
		if instructions[faultyInstructionIndex].Name == "nop" {
			switchedSomething = true
			fmt.Println("Tried switching `nop` to `jmp` statement line", faultyInstructionIndex)
			instructionsFixedPotential[faultyInstructionIndex].Name = "jmp"
		}
		if switchedSomething {
			acc, err = runProgram(instructionsFixedPotential)
		}
	}
	fmt.Println("Running fixing line: ", faultyInstructionIndex+1, ". Accumulator holds: ", acc)
}
