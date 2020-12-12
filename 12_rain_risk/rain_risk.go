package main

import (
	"fmt"
	"io/ioutil"
	"math"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type grid interface {
	getCoord() (int, int)
}

func manhattanDistance(g grid) int {
	x, y := g.getCoord()
	return int(math.Abs(float64(x)) + math.Abs(float64(y)))
}

type Position struct {
	X, Y           int
	directionIndex int
}

func (p Position) getCoord() (int, int) {
	return p.X, p.Y
}

func getSupportedDirections() []rune {
	return []rune{'N', 'E', 'S', 'W'}
}

func getInstructionValue(instruction string) int {
	value, err := strconv.Atoi(fmt.Sprint(instruction[1:len(instruction)]))
	check(err)
	return value
}

func (position *Position) processNaiveInstruction(instruction string) {
	value := getInstructionValue(instruction)
	switch instruction[0] {
	case 'N':
		position.X -= value
	case 'S':
		position.X += value
	case 'E':
		position.Y += value
	case 'W':
		position.Y -= value
	case 'F':
		absoluteInstruction := fmt.Sprintf("%s%v", string(getSupportedDirections()[position.directionIndex]), value)
		position.processNaiveInstruction(absoluteInstruction)
	case 'L':
		rEquivalent := -value
		for rEquivalent < 0 {
			rEquivalent += 360
		}
		position.processNaiveInstruction(fmt.Sprintf("R%v", rEquivalent))
	case 'R':
		position.directionIndex = (position.directionIndex + (value / 90)) % 4
	}
}

type WayPoint struct {
	N int
	E int
}

type WayPointPosition struct {
	X, Y     int
	WayPoint WayPoint
}

func (p WayPointPosition) getCoord() (int, int) {
	return p.X, p.Y
}

func (position *WayPointPosition) processWayPointInstruction(instruction string) {
	value := getInstructionValue(instruction)

	switch instruction[0] {
	case 'N':
		position.WayPoint.N += value
	case 'S':
		position.WayPoint.N -= value
	case 'E':
		position.WayPoint.E += value
	case 'W':
		position.WayPoint.E -= value
	case 'F':
		position.X += value * position.WayPoint.E
		position.Y -= value * position.WayPoint.N
	case 'L':
		rEquivalent := -value
		for rEquivalent < 0 {
			rEquivalent += 360
		}
		position.processWayPointInstruction(fmt.Sprintf("R%v", rEquivalent))
	case 'R':
		quarts := (value / 90) % 4
		switch quarts {
		case 1:
			position.WayPoint.E, position.WayPoint.N = position.WayPoint.N, -position.WayPoint.E
		case 2:
			position.WayPoint.E, position.WayPoint.N = -position.WayPoint.E, -position.WayPoint.N
		case 3:
			position.WayPoint.E, position.WayPoint.N = -position.WayPoint.N, position.WayPoint.E
		}
	}
}

func main() {
	inputFile := os.Args[1]
	input := readFile(inputFile)
	instructionStream := strings.Split(input, "\n")

	ferryPosition := Position{0, 0, 1}
	for _, instruction := range instructionStream {
		ferryPosition.processNaiveInstruction(instruction)
	}
	fmt.Println("Final ferry distance after naive instruction processing:", manhattanDistance(ferryPosition))

	wayPoint := WayPoint{1, 10}
	ferryWayPointPosition := WayPointPosition{0, 0, wayPoint}

	for _, instruction := range instructionStream {
		ferryWayPointPosition.processWayPointInstruction(instruction)
	}
	fmt.Println("Final ferry distance after waypoint instruction processing:", manhattanDistance(ferryWayPointPosition))
}
