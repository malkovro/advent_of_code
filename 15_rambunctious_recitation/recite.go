package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func processInput() (inputStream []int, desiredSpokenIndex int) {
	inputFile := os.Args[1]
	desiredSpokenIndex, _ = strconv.Atoi(os.Args[2])
	dat, err := ioutil.ReadFile(inputFile)
	check(err)
	for _, str := range strings.Split(string(dat), ",") {
		intValue, castErr := strconv.Atoi(str)
		check(castErr)
		inputStream = append(inputStream, intValue)
	}
	return
}

type Turn struct {
	Spoken int
	Count int
}

func (turn *Turn) updateSpoken(s int) {
	if turn.Count % 100000  == 0 {
		fmt.Println("Saying => ",s)
	}
	turn.Spoken = s
	turn.Count++
}

func main() {
	inputStream, desiredSpokenIndex := processInput()
	lastOccurenceIndex := make(map[int]int)

	var lastTurn = Turn{inputStream[0], 0}

	for lastTurn.Count < desiredSpokenIndex {
		var willSpeak int
		if lastTurn.Count < len(inputStream) {
			willSpeak = inputStream[lastTurn.Count]
		} else {
			value, exists := lastOccurenceIndex[lastTurn.Spoken]
			if !exists {
				willSpeak = 0
			} else {
				willSpeak = lastTurn.Count - value
			}
		}
		lastOccurenceIndex[lastTurn.Spoken] = lastTurn.Count
		lastTurn.updateSpoken(willSpeak)
	}

	fmt.Println(desiredSpokenIndex, "th number spoken is", lastTurn.Spoken)
}

