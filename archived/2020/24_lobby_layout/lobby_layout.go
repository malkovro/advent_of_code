package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func debugPrintf(str string, args ...interface{}) {
	if debug {
		fmt.Printf(str, args...)
	}
}

func debugPrintln(args ...interface{}) {
	if debug {
		fmt.Println(args...)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type LobbyFloor map[Coord]bool

func (l LobbyFloor) BlackTilesCount() (count int) {
	for _, isBlack := range l {
		if isBlack {
			count++
		}
	}
	return
}

type Coord struct {
	Q, R int
}

func (c Coord) String() string {
	return fmt.Sprintf("Q: %v, R: %v", c.Q, c.R)
}

func ParseInstructions(line string) (instructions []string) {
	buffer := ""
	for _, r := range line {
		switch r {
		case 'n':
			buffer = fmt.Sprintf("%v", string(r))
		case 's':
			buffer = fmt.Sprintf("%v", string(r))
		default:
			instructions = append(instructions, fmt.Sprintf("%v%v", string(buffer), string(r)))
			buffer = ""
		}
	}
	return
}

func (c Coord) MoveTo(line string) (dest Coord) {
	dest.Q = c.Q
	dest.R = c.R
	instructions := ParseInstructions(line)
	for _, instr := range instructions {
		switch instr {
		case "ne":
			if dest.R%2 != 0 {
				dest.Q += 1
			}
			dest.R -= 1
		case "e":
			dest.Q += 1
		case "se":
			if dest.R%2 != 0 {
				dest.Q += 1
			}
			dest.R += 1
		case "sw":
			if dest.R%2 == 0 {
				dest.Q -= 1
			}
			dest.R += 1
		case "w":
			dest.Q -= 1
		case "nw":
			if dest.R%2 == 0 {
				dest.Q -= 1
			}
			dest.R -= 1
		}
	}
	return dest
}

func (c Coord) Neighbours() []Coord {
	return []Coord{
		c.MoveTo("ne"), c.MoveTo("e"), c.MoveTo("se"),
		c.MoveTo("sw"), c.MoveTo("w"), c.MoveTo("nw"),
	}
}

func FindIndex(slice []Coord, item Coord) int {
	for i, ref := range slice {
		if ref == item {
			return i
		}
	}
	return -1
}
func AddWithoutDuplicates(coordSlice *[]Coord, coordsToAdd ...Coord) {
	missingCoords := []Coord{}
	for _, coord := range coordsToAdd {
		index := FindIndex(*coordSlice, coord)
		if index == -1 {
			missingCoords = append(missingCoords, coord)
		}
	}
	*coordSlice = append(*coordSlice, missingCoords...)
}

func TilesToCheck(blackTilesMap map[Coord]bool) []Coord {
	coordsToCheck := []Coord{}
	for coord, isBlack := range blackTilesMap {
		if isBlack {
			AddWithoutDuplicates(&coordsToCheck, coord)
			AddWithoutDuplicates(&coordsToCheck, coord.Neighbours()...)
		}
	}
	return coordsToCheck
}

func ShouldFlip(coord Coord, blackTilesMap map[Coord]bool) bool {
	adjBlackCount := 0
	for _, adj := range coord.Neighbours() {
		if blackTilesMap[adj] {
			adjBlackCount++
		}
	}
	return (blackTilesMap[coord] && (adjBlackCount == 0 || adjBlackCount > 2) || !blackTilesMap[coord] && adjBlackCount == 2)
}

var debug = false

func main() {
	defer timeTrack(time.Now(), "Solving Time")

	inputFile := flag.String("f", "this is not optional!", "the input file")
	iterations := flag.Int("it", 100, "Number of iterations")
	flag.BoolVar(&debug, "d", false, "Print the hands played")
	flag.Parse()

	fmt.Printf("==> Solving taking on file (%v):\n", *inputFile)
	dat, _ := ioutil.ReadFile(*inputFile)
	inputLines := strings.Split(string(dat), "\n")

	blackTilesMap := make(map[Coord]bool)

	refCoord := Coord{0, 0}

	fmt.Println("Number of tiles:", len(inputLines))

	for _, line := range inputLines {
		tile := refCoord.MoveTo(line)
		blackTilesMap[tile] = !blackTilesMap[tile]
	}

	for i := 0; i < *iterations; i++ {
		if debug {
			fmt.Printf("Number of black tiles after Cycle %v: %v\n", i, LobbyFloor(blackTilesMap).BlackTilesCount())
		}
		tilesToCheck := TilesToCheck(blackTilesMap)
		tilesToFlip := []Coord{}

		for _, tile := range tilesToCheck {
			if ShouldFlip(tile, blackTilesMap) {
				tilesToFlip = append(tilesToFlip, tile)
			}
		}

		debugPrintf("Flipping %v tiles\n", len(tilesToFlip))
		for _, tile := range tilesToFlip {
			blackTilesMap[tile] = !blackTilesMap[tile]
		}
	}

	fmt.Printf("Number of black tiles after %v cycles: %v\n", *iterations, LobbyFloor(blackTilesMap).BlackTilesCount())

}
