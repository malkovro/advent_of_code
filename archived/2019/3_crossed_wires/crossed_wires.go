package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

type Segment struct {
	Direction rune
	Count     int
}

type Cursor struct {
	X, Y int
}

type Distance struct {
	Man, Delay int
}

func drawSegment(canvas *map[string]int, segment Segment, cursor *Cursor, delay int) {
	switch segment.Direction {
	case 'U':
		for j, _ := range make([]int, segment.Count) {
			key := fmt.Sprintf("%v,%v", cursor.X, cursor.Y+j+1)
			(*canvas)[key] = delay + j + 1
		}
		cursor.Y += segment.Count
	case 'D':
		for j, _ := range make([]int, segment.Count) {
			key := fmt.Sprintf("%v,%v", cursor.X, cursor.Y-(j+1))
			(*canvas)[key] = delay + j
		}
		cursor.Y -= segment.Count
	case 'L':
		for i, _ := range make([]int, segment.Count) {
			key := fmt.Sprintf("%v,%v", cursor.X-(i+1), cursor.Y)
			(*canvas)[key] = delay + i + 1
		}
		cursor.X -= segment.Count
	case 'R':
		for i, _ := range make([]int, segment.Count) {
			key := fmt.Sprintf("%v,%v", cursor.X+(i+1), cursor.Y)
			(*canvas)[key] = delay + i + 1
		}
		cursor.X += segment.Count
	}
}

func checkAlongSegment(canvas map[string]int, segment Segment, cursor *Cursor, delay int) (interesections []Distance) {

	switch segment.Direction {
	case 'U':
		for j, _ := range make([]int, segment.Count) {
			key := fmt.Sprintf("%v,%v", cursor.X, cursor.Y+j+1)
			if canvas[key] != 0 {
				distance := Distance{manDistance(cursor.X, cursor.Y+j+1), delay + j + 1 + canvas[key]}
				interesections = append(interesections, distance)
			}
		}
		cursor.Y = cursor.Y + segment.Count

	case 'D':
		for j, _ := range make([]int, segment.Count) {
			key := fmt.Sprintf("%v,%v", cursor.X, cursor.Y-(j+1))
			if canvas[key] != 0 {
				distance := Distance{manDistance(cursor.X, cursor.Y-(j+1)), delay + j + 1 + canvas[key]}
				interesections = append(interesections, distance)
			}
		}
		cursor.Y = cursor.Y - segment.Count
	case 'L':
		for i, _ := range make([]int, segment.Count) {
			key := fmt.Sprintf("%v,%v", cursor.X-(i+1), cursor.Y)
			if canvas[key] != 0 {
				distance := Distance{manDistance(cursor.X-(i+1), cursor.Y), delay + i + 1 + canvas[key]}
				interesections = append(interesections, distance)
			}
		}
		cursor.X = cursor.X - segment.Count
	case 'R':
		for i, _ := range make([]int, segment.Count) {
			key := fmt.Sprintf("%v,%v", cursor.X+(i+1), cursor.Y)
			if canvas[key] != 0 {
				distance := Distance{manDistance(cursor.X+(i+1), cursor.Y), delay + i + 1 + canvas[key]}
				interesections = append(interesections, distance)
			}
		}
		cursor.X = cursor.X + segment.Count
	}
	return
}

func manDistance(x int, y int) int {
	return int(math.Abs(float64(x)) + math.Abs(float64(y)))
}

func minMan(distances []Distance) (min int) {
	for i, d := range distances {
		if i == 0 || d.Man < min {
			min = d.Man
		}
	}
	return
}

func minDelay(distances []Distance) (min int) {
	for i, d := range distances {
		if i == 0 || d.Delay < min {
			min = d.Delay
		}
	}
	return
}

func instructionToSegment(str string) Segment {
	deplacementCount, _ := strconv.Atoi(str[1:])
	return Segment{rune(str[0]), deplacementCount}
}

func main() {
	inputFile := os.Args[1]
	dat, _ := ioutil.ReadFile(inputFile)

	wires := strings.Split(string(dat), "\n")
	wireA := strings.Split(wires[0], ",")
	wireB := strings.Split(wires[1], ",")

	canvas := make(map[string]int)

	cursor := Cursor{0, 0}
	delay := 0

	for _, instruction := range wireA {
		segment := instructionToSegment(instruction)
		drawSegment(&canvas, segment, &cursor, delay)
		delay += segment.Count
	}

	cursor = Cursor{0, 0}
	intersections := []Distance{}
	wireBDelay := 0

	for _, instruction := range wireB {
		segment := instructionToSegment(instruction)
		intersections = append(intersections, checkAlongSegment(canvas, segment, &cursor, wireBDelay)...)
		wireBDelay += segment.Count
	}

	fmt.Println("Closest Manhattan Intersection", minMan(intersections))
	fmt.Println("Closest Intersection considering delay", minDelay(intersections))
}
