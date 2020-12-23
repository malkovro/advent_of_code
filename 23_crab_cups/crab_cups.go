package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
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

var debug = false
var moveCount = 0
var cupNumber = 0
var moveNumber int

var pickedUpCups = make([]int, 3)
var restCups []int

func move(cups []int, currentCupIndex *int) {
	moveCount++
	debugPrintf("-- move %v --\n", moveCount)
	debugPrintf("cups: ")
	if debug {
		for i, cup := range cups {
			if i == *currentCupIndex {
				fmt.Printf("(%v) ", cup)
			} else {
				fmt.Printf("%v ", cup)
			}
		}
		fmt.Println()
	}
	currentCup := cups[*currentCupIndex]
	pickUp(pickedUpCups, restCups, cups, *currentCupIndex)
	debugPrintf("pick up: %v\n", pickedUpCups)
	destination := getDestinationCup(pickedUpCups, currentCup)
	debugPrintf("destination: %v\n", destination)
	replaceCups(cups, restCups, destination, pickedUpCups)
	*currentCupIndex = moveCurrentCup(cups, currentCup)
	debugPrintln("")
}

func pickUp(pickedUpCups []int, restCups []int, cups []int, currentCupIndex int) {
	restIndex := 0

	pickUpTo := -1
	if currentCupIndex+3 > cupNumber-1 {
		pickUpTo = (currentCupIndex + 3) % cupNumber
	}

	for index, i := range cups {
		if pickUpTo != -1 && index <= pickUpTo {
			pickedUpCups[2-pickUpTo+index] = i
			continue
		}

		if index > currentCupIndex && index < currentCupIndex+4 {
			pickedUpCups[index-currentCupIndex-1] = i
			continue
		}
		restCups[restIndex] = i
		restIndex++
	}
}

func getDestinationCup(pickedUpCups []int, currentCup int) (destination int) {
	try := currentCup - 1
	if try == 0 {
		try = cupNumber
	}
	for {
		skip := false
		for _, cup := range pickedUpCups {
			if cup == try {
				if try > 1 {
					try--
				} else {
					try = cupNumber
				}
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		return try
	}
}

func replaceCups(cups []int, restCups []int, destination int, pickedUpCups []int) {
	index := 0
	for _, cup := range restCups {
		cups[index] = cup
		index++
		if cup == destination {
			for _, pickedUpCup := range pickedUpCups {
				cups[index] = pickedUpCup
				index++
			}
		}
	}
}

func moveCurrentCup(cups []int, currentCup int) int {
	takeNext := false
	for i, cup := range cups {
		if takeNext {
			return i
		}
		if cup == currentCup {
			takeNext = true
		}
	}
	return 0
}

func main() {
	defer timeTrack(time.Now(), "Solving Time")

	inputFile := flag.String("f", "this is not optional!", "the input file")
	v2 := flag.Bool("v2", false, "Solve the 2nd problem")
	flag.IntVar(&cupNumber, "c", 10, "Number of cups")
	flag.IntVar(&moveNumber, "m", 10, "Number of moves")
	flag.BoolVar(&debug, "d", false, "Print the hands played")
	flag.Parse()

	fmt.Printf("==> Solving for Problem v2? [%v] taking on file (%v):\n", *v2, *inputFile)
	dat, _ := ioutil.ReadFile(*inputFile)
	inputLine := string(dat)

	currentCupIndex := 0
	cups := make([]int, cupNumber)
	for i, cup := range inputLine {
		var err error
		cups[i], err = strconv.Atoi(string(cup))
		check(err)
	}

	for i := 9; i < cupNumber; i++ {
		cups[i] = i + 1
	}

	pickedUpCups = make([]int, 3)
	restCups = make([]int, cupNumber-3)

	for i := 0; i < moveNumber; i++ {
		move(cups, &currentCupIndex)
	}

	for i, cup := range cups {
		if cup == 1 {
			fmt.Printf("Cups after one multiplied: %v \n", cups[(i+1)%len(cups)]*cups[(i+2)%len(cups)])
		}
	}
	// fmt.Println("Cups", cups)
}
