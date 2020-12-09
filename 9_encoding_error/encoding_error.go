package main

import (
	"fmt"
	"io/ioutil"
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

type BufferEntry struct {
	Number       int
	PossibleSums []int
}

type Sequence struct {
	Numbers []int
	Sum     int
}

type BufferStack []BufferEntry

func (bufferStack *BufferStack) addEntry(entry int) {
	for i, bufferEntry := range *bufferStack {
		(*bufferStack)[i].PossibleSums = append(bufferEntry.PossibleSums, (bufferEntry.Number + entry))
	}
	*bufferStack = append(*bufferStack, BufferEntry{entry, []int{}})
}

func (bufferStack BufferStack) checkInputIsValid(inputValue int) bool {
	for _, bufferEntry := range bufferStack {
		for _, potentialSum := range bufferEntry.PossibleSums {
			if inputValue == potentialSum {
				return true
			}
		}
	}
	return false
}

func toIntStream(strStream []string) (intStream []int) {
	for _, str := range strStream {
		intValue, _ := strconv.Atoi(str)
		intStream = append(intStream, intValue)
	}
	return
}

func detectWrongInput(stream []int, preambleLength int) (outlier int) {
	bufferStack := BufferStack{}

	for i, inputValue := range stream {
		if i > preambleLength {
			bufferStack = bufferStack[1:]
			isValid := bufferStack.checkInputIsValid(inputValue)

			if !isValid {
				// fmt.Println(bufferStack)
				outlier = inputValue
				return
			}
		}
		bufferStack.addEntry(inputValue)
	}
	return
}

func detectSequenceSummingTo(expectedSum int, stream []int) Sequence {
	sequences := []Sequence{}
	for _, inputValue := range stream {
		nextSequences := []Sequence{}
		for i, _ := range make([]int, len(sequences)) {
			sequence := sequences[i]

			if sequence.Sum+inputValue > expectedSum {
				continue
			}
			sequence.Numbers = append(sequence.Numbers, inputValue)
			sequence.Sum = sequence.Sum + inputValue

			if sequence.Sum == expectedSum {
				return sequence
			}

			nextSequences = append(nextSequences, sequence)
		}
		sequences = append(nextSequences, Sequence{[]int{inputValue}, inputValue})
	}
	return Sequence{}
}

func (sequence Sequence) min() (min int) {
	for i, item := range sequence.Numbers {
		if i == 0 || item < min {
			min = item
		}
	}
	return
}

func (sequence Sequence) max() (max int) {
	for i, item := range sequence.Numbers {
		if i == 0 || item > max {
			max = item
		}
	}
	return
}

func main() {
	inputFile := os.Args[1]
	preambleLength, _ := strconv.Atoi(os.Args[2])
	input := readFile(inputFile)
	stream := toIntStream(strings.Split(input, "\n"))

	detectedOutlier := detectWrongInput(stream, preambleLength)
	fmt.Println("Invalid Number detected: ", detectedOutlier)

	sequence := detectSequenceSummingTo(detectedOutlier, stream)

	fmt.Println("Sequence summing to it:", sequence.min()+sequence.max())
}
