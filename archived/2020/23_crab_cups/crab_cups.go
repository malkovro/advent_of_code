package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
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

var debug = false
var moveCount = 0
var cupNumber = 0
var moveNumber int

type CupLinkedList map[int]int

func (cs CupLinkedList) OrderedCup() []int {
	orderedCup := make([]int, cupNumber)
	orderedCup[0] = firstCup
	for i := 1; i < cupNumber; i++ {
		orderedCup[i] = cs[orderedCup[i-1]]
	}

	return orderedCup
}

func (cs CupLinkedList) ToString(highlighted int) string {
	cupStrArray := make([]string, cupNumber)
	for i, c := range cs.OrderedCup() {
		if c == highlighted {
			cupStrArray[i] = fmt.Sprintf("(%v)", c)
		} else {
			cupStrArray[i] = fmt.Sprintf("%v", c)
		}
	}
	return strings.Join(cupStrArray, " ")
}

func move(cups *CupLinkedList, currentCup *int) {
	moveCount++
	// debugPrintf("-- move %v --\n", moveCount)
	// debugPrintf("cups: ")
	// debugPrintln(cups.ToString(*currentCup))
	takes3(*cups, *currentCup)
	// debugPrintf("pick up: %v\n", pickedUpCups)
	destination := getDestinationCup(*currentCup, &pickedUpCups)
	// debugPrintf("destination: %v\n", destination)
	move3Cups(*cups, &pickedUpCups, *currentCup, destination)
	*currentCup = (*cups)[*currentCup]
	// debugPrintln(cups.ToString(*currentCup))
	debugPrintln("")
}

func takes3(cups CupLinkedList, currentCup int) {
	for i := 0; i < 3; i++ {
		pickedUpCups[i] = cups[currentCup]
		currentCup = pickedUpCups[i]
	}
}

func getDestinationCup(currentCup int, pickedUpCups *[]int) (destination int) {
	try := currentCup - 1
	if try == 0 {
		try = cupNumber
	}
	for {
		skip := false
		for _, cup := range *pickedUpCups {
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

var pickedUpCups = make([]int, 3)

func move3Cups(cups CupLinkedList, pickedUpCups *[]int, from int, to int) {
	cups[from] = cups[(*pickedUpCups)[2]]
	cups[to], cups[(*pickedUpCups)[2]] = (*pickedUpCups)[0], cups[to]
}

var firstCup int

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

	cupLList := make(CupLinkedList)

	var lastCup int
	for i, cupStr := range inputLine {
		cup, err := strconv.Atoi(string(cupStr))
		check(err)
		if i != 0 {
			cupLList[lastCup] = cup
		} else {
			firstCup = cup
		}
		lastCup = cup
	}

	if *v2 {
		for i := 10; i < cupNumber+1; i++ {
			cupLList[lastCup] = i
			lastCup = i
		}
	}

	cupLList[lastCup] = firstCup
	currentCup := firstCup

	for i := 0; i < moveNumber; i++ {
		move(&cupLList, &currentCup)
	}

	if *v2 {
		fmt.Printf("Cups after one multiplied: %v \n", cupLList[1]*cupLList[cupLList[1]])
	}
}
