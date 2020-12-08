package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func runInstruction(index int, instructions *[]int) {
	instruction := (*instructions)[index]

	destinationIndex := (*instructions)[index+3]
	aIndex := (*instructions)[index+1]
	bIndex := (*instructions)[index+2]

	switch instruction {
	case 1:
		(*instructions)[destinationIndex] = (*instructions)[aIndex] + (*instructions)[bIndex]
	case 2:
		(*instructions)[destinationIndex] = (*instructions)[aIndex] * (*instructions)[bIndex]
	default:
		panic("Invalid instruction " + fmt.Sprint(instruction))
	}
}

func setState(instructions *[]int, noun int, verb int) {
	(*instructions)[1] = noun
	(*instructions)[2] = verb
}

func cloneInstructionsWithBase(baseInstructions []int, noun int, verb int) []int {
	instructions := make([]int, len(baseInstructions))
	copy(instructions, baseInstructions)
	setState(&instructions, noun, verb)
	return instructions
}

func runInstructionsWithBase(baseInstructions []int, noun int, verb int) []int {
	instructions := cloneInstructionsWithBase(baseInstructions, noun, verb)
	currentIndex := 0
	for instructions[currentIndex] != 99 {
		runInstruction(currentIndex, &instructions)
		currentIndex += 4
	}
	return instructions
}

func main() {
	desiredOutput := 19690720
	inputFile := os.Args[1]
	dat, err := ioutil.ReadFile(inputFile)
	check(err)

	strInstructions := strings.Split(strings.Replace(string(dat), "\n", "", -1), ",")

	baseInstructions := []int{}
	for _, instructionStr := range strInstructions {
		instruction, err := strconv.Atoi(instructionStr)
		check(err)
		baseInstructions = append(baseInstructions, instruction)
	}

	instructions := runInstructionsWithBase(baseInstructions, 12, 2)
	fmt.Println("Program Halts with 0 indexed value: ", fmt.Sprint(instructions[0]))

	for noun, _ := range make([]int, 100) {
		for verb, _ := range make([]int, 100) {
			instructions := runInstructionsWithBase(baseInstructions, noun, verb)

			if instructions[0] == desiredOutput {
				fmt.Println("Found the right noun and verb:", noun*100+verb)
				return
			}
		}
	}
}
