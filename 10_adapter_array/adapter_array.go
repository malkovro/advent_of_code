package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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

func toIntStream(strStream []string) (intStream []int) {
	for _, str := range strStream {
		intValue, _ := strconv.Atoi(str)
		intStream = append(intStream, intValue)
	}
	return
}

type AdapterSet struct {
	LastJoltage int
}

func compatibility(adapter int, plug int) (int, error) {
	switch adapter - plug {
	case 1:
		return 1, nil
	case 2:
		return 2, nil
	case 3:
		return 3, nil
	default:
		str := fmt.Sprintf("We cannot plug %d into %d", adapter, plug)
		return -1, errors.New(str)
	}
}

func main() {
	inputFile := os.Args[1]
	input := readFile(inputFile)
	adapters := toIntStream(strings.Split(input, "\n"))

	sort.Slice(adapters, func(i, j int) bool {
		return adapters[i] < adapters[j]
	})

	joltDifferenceMap := make(map[int]int)
	previousAdapter := 0
	for _, adapter := range adapters {
		joltDifference, err := compatibility(adapter, previousAdapter)
		if err != nil {
			panic(err)
		}
		joltDifferenceMap[joltDifference] += 1
		previousAdapter = adapter
	}

	// Add the device built in jolt jump:
	joltDifferenceMap[3] += 1
	fmt.Println("JoltDifference:", joltDifferenceMap)
	fmt.Println("1-jolt x 3-jolt jumps = ", joltDifferenceMap[1]*joltDifferenceMap[3])

	// Holds for each available joltage the number of arrangements possible
	arrangements := make(map[int]int)

	// Add initial output plug:
	arrangements[0] = 1
	lastAdapterJoltage := adapters[len(adapters)-1]

	for i, adapter := range adapters {
		droppable := i < len(adapters)-1 // Last adapter is not droppable as device joltage is +3

		for _, joltage := range getJoltages(arrangements) {
			joltDifference, err := compatibility(adapter, joltage)
			adapterIsCompatible := err == nil

			numberOfPathToJoltage := arrangements[joltage]

			if !adapterIsCompatible {
				delete(arrangements, joltage)
				continue
			}

			if !droppable || joltDifference == 3 {
				// The current adapter cannot be dropped and we cannot branch from here
				delete(arrangements, joltage)
			}

			arrangements[adapter] += numberOfPathToJoltage
		}
	}

	fmt.Println("We appear to have ", arrangements[lastAdapterJoltage], "possible combination 😱")
}

func getJoltages(mymap map[int]int) []int {
	keys := make([]int, len(mymap))

	i := 0
	for k := range mymap {
		keys[i] = k
		i++
	}
	return keys
}
