package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

func homogenizeGroupAnswer(passport string) string {
	return strings.Replace(passport, "\n", "", -1)
}

func unique(runeSlice []rune) []rune {
	keys := make(map[rune]bool)
	list := []rune{}
	for _, entry := range runeSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func repeatedNTime(runeSlice []rune, n int) []rune {
	keys := make(map[rune]int)
	list := []rune{}

	if n == 1 {
		return unique(runeSlice)
	}

	for _, entry := range runeSlice {
		if count, exists := keys[entry]; !exists {
			keys[entry] = 1
		} else if count < (n - 1) {
			keys[entry] = count + 1
		} else {
			keys[entry] = count + 1
			list = append(list, entry)
		}
	}
	return list
}

func distinctRuneCount(answers string) int {
	return len(unique([]rune(answers)))
}

func unanimouslyAnsweredCount(answers string, n int) int {
	return len(repeatedNTime([]rune(answers), n))
}

func main() {
	inputFile := os.Args[1]
	answers := readFile(inputFile)
	groupedAnswers := strings.Split(answers, "\n\n")

	sumAnyAnsweredQuestions := 0
	sumAllAnsweredQuestions := 0

	for _, groupAnswers := range groupedAnswers {
		peopleInGroup := strings.Count(groupAnswers, "\n") + 1
		homogenizedGroupAnswers := homogenizeGroupAnswer(groupAnswers)

		sumAnyAnsweredQuestions += distinctRuneCount(homogenizedGroupAnswers)
		sumAllAnsweredQuestions += unanimouslyAnsweredCount(homogenizedGroupAnswers, peopleInGroup)
	}
	fmt.Println("Total sum of yessed questions: ", sumAnyAnsweredQuestions)
	fmt.Println("Total sum of unanimously yessed questions: ", sumAllAnsweredQuestions)
}
