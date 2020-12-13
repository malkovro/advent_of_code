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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type BusDeparture struct {
	Id, Wait int
}

func toFreqStream(busFrequenciesStr []string) (busFreqStream []int) {
	for _, busFrequency := range busFrequenciesStr {
		busId, err := strconv.Atoi(busFrequency)
		if err != nil {
			busFreqStream = append(busFreqStream, -1)
		} else {
			busFreqStream = append(busFreqStream, busId)
		}
	}
	return
}

func extendedGcd(a int, b int) (int, int, int) {
	old_r, r := a, b
	old_s, s := 1, 0
	old_t, t := 0, 1

	for r != 0 {
		quotient := old_r / r
		old_r, r = r, old_r-quotient*r
		old_s, s = s, old_s-quotient*s
		old_t, t = t, old_t-quotient*t
	}

	// fmt.Println("Bézout coefficients:", old_s, old_t)
	// fmt.Println("greatest common divisor:", old_r)
	// fmt.Println("quotients by the gcd:", t, s)

	return old_s, old_t, old_r
}

func main() {
	inputFile := os.Args[1]
	input := readFile(inputFile)
	inputData := strings.Split(input, "\n")
	startTime, _ := strconv.Atoi(inputData[0])
	busFrequencies := toFreqStream(strings.Split(inputData[1], ","))

	nearestDeparture := BusDeparture{}
	for i, busId := range busFrequencies {
		if busId == -1 {
			continue
		}

		nextBusDeparture := busId - startTime%busId
		// fmt.Println("Bus", busId, "would leave in ", nextBusDeparture)
		if i == 0 || nextBusDeparture < nearestDeparture.Wait {
			nearestDeparture = BusDeparture{busId, nextBusDeparture}
		}
	}
	fmt.Println("Bus", nearestDeparture.Id, " would leave in ", nearestDeparture.Wait)
	fmt.Println("Bus Id of nearest departure x waiting time =", nearestDeparture.Id*nearestDeparture.Wait)

	multipliedFrequencies := 1
	for _, busId := range busFrequencies {
		if busId != -1 {
			multipliedFrequencies = multipliedFrequencies * busId
		}
	}
	fmt.Println("N:", multipliedFrequencies)
	x := 0
	for index, busId := range busFrequencies {
		if busId == -1 {
			continue
		}

		restMultiplier := multipliedFrequencies / busId
		_, s, _ := extendedGcd(busId, restMultiplier)

		fmt.Println("Bus Id:", busId, "| Ai:", busId-index, "| Si:", s, "| N/ni:", restMultiplier)
		x += (busId - index) * s * restMultiplier
	}

	res := x % multipliedFrequencies
	for res < 0 {
		res += multipliedFrequencies
	}
	fmt.Println("Thanks maths, the magic timestamp is: ", res)
}
