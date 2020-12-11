package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

type Coord struct {
	X, Y int
}

type WaitingArea struct {
	SeatMap map[Coord]bool
	X, Y    int
}

func (waitingArea WaitingArea) get(coord Coord) (bool, bool) {
	v, e := waitingArea.SeatMap[coord]
	return v, e
}

func (waitingArea *WaitingArea) set(coord Coord, value bool) {
	(*waitingArea).SeatMap[coord] = value
}

func buildSeatMap(strStream []string) WaitingArea {
	seatMap := make(map[Coord]bool)
	for j, str := range strStream {
		for i, seatDefinition := range str {
			if seatDefinition == '.' {
				continue
			}
			seatMap[Coord{i, j}] = false
		}
	}
	return WaitingArea{seatMap, len(strStream[0]), len(strStream)}
}

func hasRightNeighbour(coord Coord, seatMap WaitingArea) bool {
	for i, _ := range make([]int, seatMap.X-coord.X) {
		v, exists := seatMap.get(Coord{coord.X + i + 1, coord.Y})
		if exists {
			return v
		}
	}
	return false
}

func hasLeftNeighbour(coord Coord, seatMap WaitingArea) bool {
	for i, _ := range make([]int, coord.X) {
		v, exists := seatMap.get(Coord{coord.X - (i + 1), coord.Y})
		if exists {
			return v
		}
	}
	return false
}

func hasTopNeighbour(coord Coord, seatMap WaitingArea) bool {
	for i, _ := range make([]int, coord.Y) {
		v, exists := seatMap.get(Coord{coord.X, coord.Y - (i + 1)})
		if exists {
			return v
		}
	}
	return false
}

func hasBottomNeighbour(coord Coord, seatMap WaitingArea) bool {
	for i, _ := range make([]int, seatMap.Y-coord.Y) {
		v, exists := seatMap.get(Coord{coord.X, coord.Y + (i + 1)})
		if exists {
			return v
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func hasBottomRightNeighbour(coord Coord, seatMap WaitingArea) bool {
	maxDirection := min(seatMap.Y-coord.Y, seatMap.X-coord.X)
	for i, _ := range make([]int, maxDirection) {
		v, exists := seatMap.get(Coord{coord.X + (i + 1), coord.Y + (i + 1)})
		if exists {
			return v
		}
	}
	return false
}

func hasBottomLeftNeighbour(coord Coord, seatMap WaitingArea) bool {
	maxDirection := min(seatMap.Y-coord.Y, coord.X)
	for i, _ := range make([]int, maxDirection) {
		v, exists := seatMap.get(Coord{coord.X - (i + 1), coord.Y + (i + 1)})
		if exists {
			return v
		}
	}
	return false
}

func hasTopLeftNeighbour(coord Coord, seatMap WaitingArea) bool {
	maxDirection := min(coord.Y, coord.X)
	for i, _ := range make([]int, maxDirection) {
		v, exists := seatMap.get(Coord{coord.X - (i + 1), coord.Y - (i + 1)})
		if exists {
			return v
		}
	}
	return false
}

func hasTopRightNeighbour(coord Coord, seatMap WaitingArea) bool {
	maxDirection := min(coord.Y, seatMap.X-coord.X)
	for i, _ := range make([]int, maxDirection) {
		v, exists := seatMap.get(Coord{coord.X + (i + 1), coord.Y - (i + 1)})
		if exists {
			return v
		}
	}
	return false
}

type NeighbourDetector = func(Coord, WaitingArea) bool

func directiveNeighboursCount(coord Coord, seatMap WaitingArea) (count int) {
	neighbourFns := []NeighbourDetector{
		hasTopLeftNeighbour, hasTopNeighbour,
		hasTopRightNeighbour, hasLeftNeighbour,
		hasRightNeighbour, hasBottomLeftNeighbour,
		hasBottomRightNeighbour, hasBottomNeighbour,
	}
	for _, neighbourFn := range neighbourFns {
		if neighbourFn(coord, seatMap) {
			count++
		}
	}
	return
}

func directNeighboursCount(coord Coord, seatMap WaitingArea) (count int) {
	neighbours := []Coord{
		Coord{coord.X + 1, coord.Y + 1},
		Coord{coord.X + 1, coord.Y},
		Coord{coord.X + 1, coord.Y - 1},
		Coord{coord.X, coord.Y + 1},
		Coord{coord.X, coord.Y - 1},
		Coord{coord.X - 1, coord.Y + 1},
		Coord{coord.X - 1, coord.Y - 1},
		Coord{coord.X - 1, coord.Y},
	}
	for _, coord := range neighbours {
		if v, _ := seatMap.get(coord); v {
			count++
		}
	}
	return
}

func (area *WaitingArea) tick(tolerance int) bool {
	isStable := true
	movingCoords := []Coord{}
	for coord, seatOccupied := range area.SeatMap {
		neighboursCount := directNeighboursCount(coord, *area)
		if !seatOccupied && neighboursCount == 0 {
			movingCoords = append(movingCoords, coord)
		} else if seatOccupied && neighboursCount >= tolerance {
			movingCoords = append(movingCoords, coord)
		}
	}

	for _, coord := range movingCoords {
		isStable = false
		v, _ := area.get(coord)
		area.set(coord, !v)
	}
	return isStable
}

func (area *WaitingArea) directiveTick(tolerance int) bool {
	isStable := true
	movingCoords := []Coord{}
	for coord, seatOccupied := range area.SeatMap {

		neighboursCount := directiveNeighboursCount(coord, *area)
		if !seatOccupied && neighboursCount == 0 {
			movingCoords = append(movingCoords, coord)
		} else if seatOccupied && neighboursCount >= tolerance {
			movingCoords = append(movingCoords, coord)
		}
	}

	for _, coord := range movingCoords {
		isStable = false
		v, _ := area.get(coord)
		area.set(coord, !v)
	}
	return isStable
}

func (area WaitingArea) countOccupiedSeat() (count int) {
	for _, seatOccupied := range area.SeatMap {
		if seatOccupied {
			count++
		}
	}
	return
}

func (area WaitingArea) print() {
	fmt.Print("┌")
	for range make([]rune, area.X) {
		fmt.Print("-")
	}
	fmt.Println("┐")

	for j, _ := range make([]int, area.Y) {
		fmt.Print("|")
		for i, _ := range make([]int, area.X) {
			v, exists := area.SeatMap[Coord{i, j}]
			if exists {
				if v {
					fmt.Print("#")
				} else {
					fmt.Print("L")
				}
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("|")
		fmt.Println()
	}
	fmt.Print("└")
	for range make([]rune, area.X) {
		fmt.Print("-")
	}
	fmt.Println("┘")
}

func main() {
	inputFile := os.Args[1]
	input := readFile(inputFile)
	seatMap := buildSeatMap(strings.Split(input, "\n"))
	seatMapPb1 := buildSeatMap(strings.Split(input, "\n"))

	isStable := false

	// Problem 1
	for !isStable {
		isStable = seatMapPb1.tick(4)
	}
	fmt.Println("Total number of occupied Seats (direct neighbour) ", seatMapPb1.countOccupiedSeat())

	// Problem 2
	isStable = false
	for !isStable {
		isStable = seatMap.directiveTick(5)
		seatMap.print()
	}

	fmt.Println("Total number of occupied Seats ", seatMap.countOccupiedSeat())

}
