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

func getSupportedDirections() []rune {
	return []rune{'N', 'E', 'S', 'W'}
}

type Position struct {
	X, Y           int
	directionIndex int
}

type WayPoint struct {
	N int
	E int
}

type WayPointPosition struct {
	X, Y     int
	WayPoint WayPoint
}

func (position Position) manhattanDistance() int {
	return int(math.Abs(float64(position.X)) + math.Abs(float64(position.Y)))
}
func (position WayPointPosition) manhattanDistance() int {
	return int(math.Abs(float64(position.X)) + math.Abs(float64(position.Y)))
}

func (position *WayPointPosition) processWayPointInstruction(instruction string) {
	value, err := strconv.Atoi(fmt.Sprint(instruction[1:len(instruction)]))
	check(err)

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
		rEquivalentInstruction := fmt.Sprintf("R%v", rEquivalent)
		position.processWayPointInstruction(rEquivalentInstruction)
	case 'R':
		quarts := (value / 90) % 4
		switch quarts {
		case 0:
			return
		case 1:
			position.WayPoint.E, position.WayPoint.N = position.WayPoint.N, -position.WayPoint.E
		case 2:
			position.WayPoint.E, position.WayPoint.N = -position.WayPoint.E, -position.WayPoint.N
		case 3:
			position.WayPoint.E, position.WayPoint.N = -position.WayPoint.N, position.WayPoint.E
		}
	}
}

func (position *Position) processNaiveInstruction(instruction string) {
	value, err := strconv.Atoi(fmt.Sprint(instruction[1:len(instruction)]))
	check(err)

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
		position.directionIndex = (position.directionIndex - (value / 90)) % 4
		for position.directionIndex < 0 {
			position.directionIndex += 4
		}
	case 'R':
		position.directionIndex = (position.directionIndex + (value / 90)) % 4
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
	fmt.Println("Finale Ferry position after naive instruction processing ", ferryPosition.manhattanDistance())

	wayPoint := WayPoint{1, 10}
	ferryWayPointPosition := WayPointPosition{0, 0, wayPoint}

	for _, instruction := range instructionStream {
		ferryWayPointPosition.processWayPointInstruction(instruction)
	}
	fmt.Println("Finale Ferry position after Waypoint instruction processing ", ferryWayPointPosition.manhattanDistance())

}
