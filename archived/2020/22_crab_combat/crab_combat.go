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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func play(deckA *[]int, deckB *[]int) bool {
	deckALen := len(*deckA)
	deckBLen := len(*deckB)
	if deckALen == 0 || deckBLen == 0 {
		return false
	}
	if (*deckA)[0] > (*deckB)[0] {
		*deckA = append((*deckA)[1:deckALen], []int{(*deckA)[0], (*deckB)[0]}...)
		*deckB = (*deckB)[1:deckBLen]
	} else {
		*deckB = append((*deckB)[1:deckBLen], []int{(*deckB)[0], (*deckA)[0]}...)
		*deckA = (*deckA)[1:deckALen]
	}
	return true
}

func score(deck []int) int {
	score := 0
	multiplier := 1
	deckLength := len(deck)
	for i, _ := range make([]int, deckLength) {
		score += multiplier * deck[deckLength-i-1]
		multiplier++
	}
	return score
}

var totalGameNumber = 0
var debug = false

func main() {
	defer timeTrack(time.Now(), "Solving Time")
	inputFile := flag.String("f", "this is not optional!", "the input file")
	v2 := flag.Bool("v2", false, "Solve the 2nd problem")
	flag.BoolVar(&debug, "d", false, "Print the hands played")
	flag.Parse()

	fmt.Printf("==> Solving for Problem v2? [%v] taking on file (%v):\n", *v2, *inputFile)
	dat, _ := ioutil.ReadFile(*inputFile)
	lines := strings.Split(string(dat), "\n")

	myDeck := []int{}
	crabDeck := []int{}

	currentDeck := &myDeck
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		if strings.HasPrefix(line, "Player") {
			currentDeck = &crabDeck
			continue
		}
		card, err := strconv.Atoi(line)
		check(err)
		*currentDeck = append(*currentDeck, card)
	}

	if *v2 {
		playRecursice(&myDeck, &crabDeck)
	} else {
		for play(&myDeck, &crabDeck) {
		}
	}

	fmt.Println("Deck's scores", score(myDeck), score(crabDeck))
}

func deckSignature(deck []int) int {
	baseMultiplier := 54
	signature := 0
	for _, card := range deck {
		signature += card * baseMultiplier
		baseMultiplier *= 54
	}
	return signature
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

func playRecursice(deckA *[]int, deckB *[]int) string {
	totalGameNumber++
	gameNumber := totalGameNumber
	playedSets := make(map[[2]int]bool)
	debugPrintf("\n=== Game %v ===\n\n", gameNumber)
	round := 0
	var winner string
	for winner == "" {
		round++
		debugPrintf("\n-- Round %v (Game %v) --\n", round, gameNumber)
		debugPrintln("Player 1's deck:", *deckA)
		debugPrintln("Player 2's deck:", nil, *deckB)

		if len(*deckA) == 0 {
			return "crab"
		}
		if len(*deckB) == 0 {
			return "me"
		}
		setSignature := [2]int{deckSignature(*deckA), deckSignature(*deckB)}

		if _, exists := playedSets[setSignature]; exists {
			debugPrintln("Infinite loop => I'll take the win")
			return "me"
		}
		playedSets[setSignature] = true
		playRound(deckA, deckB, round, gameNumber)

	}
	return winner
}

func playRound(deckA *[]int, deckB *[]int, round int, gameNumber int) {
	myDeckLen := len(*deckA)
	crabDeckLen := len(*deckB)
	myCard := (*deckA)[0]
	crabCard := (*deckB)[0]

	var roundWinner string
	if myCard > myDeckLen-1 || crabCard > crabDeckLen-1 {
		if myCard > crabCard {
			roundWinner = "me"
		} else {
			roundWinner = "crab"
		}
	} else {
		debugPrintln("Playing a sub-game to determine the winner...")
		subDeckA := make([]int, myCard)
		subDeckB := make([]int, crabCard)
		copy(subDeckA, (*deckA)[1:myCard+1])
		copy(subDeckB, (*deckB)[1:crabCard+1])
		roundWinner = playRecursice(&subDeckA, &subDeckB)
		debugPrintf("...anyway, back to game %v.\n", gameNumber)
	}

	if roundWinner == "me" {
		*deckA = append((*deckA)[1:myDeckLen], []int{(*deckA)[0], (*deckB)[0]}...)
		*deckB = (*deckB)[1:crabDeckLen]
		debugPrintf("Me wins round %v of game %v!\n", round, gameNumber)
	} else {
		*deckB = append((*deckB)[1:crabDeckLen], []int{(*deckB)[0], (*deckA)[0]}...)
		*deckA = (*deckA)[1:myDeckLen]
		debugPrintf("Crab wins round %v of game %v!\n", round, gameNumber)
	}

}
