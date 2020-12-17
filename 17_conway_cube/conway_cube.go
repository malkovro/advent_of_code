package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Slice struct {
	x, y, z int
	plan    [][]bool
}

type Coord struct {
	X, Y, Z, W int
}

type Grid struct {
	min, max Coord
	Plan     map[Coord]bool
}

func (g Grid) neighbourCount(coord Coord) (count int) {
	neighboursOffsetSlice := make([]int, 3)
	comparingCoord := Coord{}
	for w, _ := range neighboursOffsetSlice {
		comparingCoord.W = coord.W - 1 + w
		for z, _ := range neighboursOffsetSlice {
			comparingCoord.Z = coord.Z - 1 + z
			for y, _ := range neighboursOffsetSlice {
				comparingCoord.Y = coord.Y - 1 + y

				for x, _ := range neighboursOffsetSlice {
					comparingCoord.X = coord.X - 1 + x
					value, _ := g.Plan[comparingCoord]
					if value && comparingCoord != coord {
						count++
					}
				}
			}
		}
	}
	return
}

func processInput() Grid {
	var y, x, z, w int
	inputFile := os.Args[1]
	dat, _ := ioutil.ReadFile(inputFile)

	splittedInputs := strings.Split(string(dat), "\n")
	grid := Grid{Coord{0, 0, 0, 0}, Coord{0, 0, 0, 0}, make(map[Coord]bool)}

	z, w = 0, 0
	y = len(splittedInputs)
	plan := make([][]bool, y)
	for j, line := range splittedInputs {
		x = len(line)
		slice := make([]bool, x)
		plan[j] = slice
		for i, state := range line {
			switch state {
			case '#':
				coord := Coord{i, j, z, w}
				grid.Plan[coord] = true
				grid.UpdateBounds(coord)
			}
		}
	}

	return grid
}

func (g *Grid) UpdateBounds(coord Coord) {
	if coord.X < g.min.X {
		(*g).min.X = coord.X
	}
	if coord.Y < g.min.Y {
		(*g).min.Y = coord.Y
	}
	if coord.Z < g.min.Z {
		(*g).min.Z = coord.Z
	}
	if coord.W < g.min.W {
		(*g).min.W = coord.W
	}
	if coord.X > g.max.X {
		(*g).max.X = coord.X
	}
	if coord.Y > g.max.Y {
		(*g).max.Y = coord.Y
	}
	if coord.Z > g.max.Z {
		(*g).max.Z = coord.Z
	}
	if coord.W > g.max.W {
		(*g).max.W = coord.W
	}
}

type Plot []Coord

func (p Plot) ToGrid() Grid {
	grid := Grid{Coord{0, 0, 0, 0}, Coord{0, 0, 0, 0}, make(map[Coord]bool)}
	for _, coord := range p {
		grid.Plan[coord] = true
		grid.UpdateBounds(coord)
	}
	return grid
}

func runCycle(grid Grid, dim4 bool) []Coord {
	nextActiveCoords := []Coord{}

	wSlice := make([]int, grid.max.W-grid.min.W+3)
	if !dim4 {
		wSlice = []int{0}
	} else {
		for l, _ := range wSlice {
			wSlice[l] = grid.min.W - 1 + l
		}
	}

	zSlice := make([]int, grid.max.Z-grid.min.Z+3)
	ySlice := make([]int, grid.max.Y-grid.min.Y+3)
	xSlice := make([]int, grid.max.X-grid.min.X+3)
	consideredCoord := Coord{}
	for _, w := range wSlice {
		consideredCoord.W = w
		for k, _ := range zSlice {
			consideredCoord.Z = grid.min.Z - 1 + k
			for j, _ := range ySlice {
				consideredCoord.Y = grid.min.Y - 1 + j
				for i, _ := range xSlice {
					consideredCoord.X = grid.min.X - 1 + i
					neighboursCount := grid.neighbourCount(consideredCoord)
					value, _ := grid.Plan[consideredCoord]

					if value && neighboursCount >= 2 && neighboursCount <= 3 {
						nextActiveCoords = append(nextActiveCoords, consideredCoord)
						continue
					}

					if !value && neighboursCount == 3 {
						nextActiveCoords = append(nextActiveCoords, consideredCoord)
						continue
					}
				}
			}
		}
	}
	return nextActiveCoords
}

func (grid Grid) Print() {
	wSlice := make([]int, grid.max.W-grid.min.W+1)
	zSlice := make([]int, grid.max.Z-grid.min.Z+1)
	ySlice := make([]int, grid.max.Y-grid.min.Y+1)
	xSlice := make([]int, grid.max.X-grid.min.X+1)

	coord := Coord{}
	for l, _ := range wSlice {
		coord.W = grid.min.W + l
		for k, _ := range zSlice {
			coord.Z = grid.min.Z + k
			fmt.Printf("\n\nz=%v, w=%v\n", coord.Z, coord.W)
			for j, _ := range ySlice {
				coord.Y = grid.min.Y + j
				for i, _ := range xSlice {
					coord.X = grid.min.X + i
					value, _ := grid.Plan[coord]
					if value {
						fmt.Print("#")
					} else {
						fmt.Print(".")
					}
				}
				fmt.Print("\n")
			}
		}
	}
}

func main() {
	grid := processInput()

	var coords []Coord
	for _ = range make([]int, 6) {
		coords = runCycle(grid, false)
		grid = Plot(coords).ToGrid()
		// grid.Print()
	}
	fmt.Printf("Number of active cubes in 3D after 6 cycle(s): %v\n", len(coords))

	grid = processInput()
	for _ = range make([]int, 6) {
		coords = runCycle(grid, true)
		grid = Plot(coords).ToGrid()
		// grid.Print()
	}
	fmt.Printf("Number of active cubes in 4D after 6 cycle(s): %v\n", len(coords))
}
