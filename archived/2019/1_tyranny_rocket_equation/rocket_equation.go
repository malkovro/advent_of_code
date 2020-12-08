package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func computeFuelRequired(weight int) int {
	return weight/3 - 2
}

func computeFuelRequiredWithSelf(weight int) int {
	fuelRequirement := computeFuelRequired(weight)

	indirectFuelRequirement := computeFuelRequired(fuelRequirement)
	for indirectFuelRequirement > 0 {
		fuelRequirement += indirectFuelRequirement
		indirectFuelRequirement = computeFuelRequired(indirectFuelRequirement)
	}
	return fuelRequirement
}

func main() {
	inputFile := os.Args[1]
	dat, err := ioutil.ReadFile(inputFile)
	check(err)
	requiredFuelTotal := 0
	requiredFuelTotalWithSelf := 0

	for _, weightStr := range strings.Split(string(dat), "\n") {
		weight, err := strconv.Atoi(weightStr)
		check(err)
		requiredFuelTotal += computeFuelRequired(weight)
		requiredFuelTotalWithSelf += computeFuelRequiredWithSelf(weight)
	}

	fmt.Println("Required total fuel is: ", requiredFuelTotal)
	fmt.Println("Required total fuel (counting fuel needed for fuel 🙃) is: ", requiredFuelTotalWithSelf)
}
