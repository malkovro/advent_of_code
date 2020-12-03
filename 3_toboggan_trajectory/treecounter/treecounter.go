package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

type Position struct {
  X,Y int
}

func replaceAtIndex(in string, r rune, i int) string {
    out := []rune(in)
    out[i] = r
    return string(out)
}

func treeEncounters(layers []string, vInc int, hInc int) (treeCount int) {
    treeCount = 0
    mapPatternSize := len(layers[0])
    pos := Position{0, 0}

    for index := 0; index < len(layers); index+=vInc {
      element := layers[index]

      if len(element) != mapPatternSize {
	break
      }
      if patternIndex := pos.Y % mapPatternSize; element[patternIndex] == '#' {
        treeCount++
	// fmt.Println(replaceAtIndex(element, 'O', patternIndex))
      } else {
        // fmt.Println(replaceAtIndex(element, 'X', patternIndex))
      }
      pos = Position{index, pos.Y+hInc}
    }
    return
}

func main() {
    inputFile := os.Args[1]
    fmt.Printf("Reading map from: %v \n", inputFile)

    data, err := ioutil.ReadFile(inputFile)
    if err != nil {
        fmt.Println("File reading error", err)
        return
    }

    layers := strings.Split(string(data), "\n")

    slop1_1 := treeEncounters(layers, 1, 1)
    slop1_3 := treeEncounters(layers, 1, 3)
    slop1_5 := treeEncounters(layers, 1, 5)
    slop1_7 := treeEncounters(layers, 1, 7)
    slop2_1 := treeEncounters(layers, 2, 1)

    fmt.Println("Simple Slope (right 3 down 1) :", slop1_3)
    fmt.Println("Slope Multiplied: ", slop1_1 * slop1_3 * slop1_5 * slop1_7 * slop2_1)
}
