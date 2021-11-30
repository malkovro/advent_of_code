package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Operation struct {
	operation string
	buffer    int
}

func (o Operation) Compute(a, b int) int {
	switch o.operation {
	case "add":
		return a + b
	default:
		return a * b
	}
}

type Stack []Operation

func (s Stack) Last() *Operation {
	return &s[len(s)-1]
}

func (s *Stack) addOperation(operation string) {
	*s = append(*s, Operation{operation, 0})
}

func (s *Stack) popOperation() {
	index := len(*s) - 1
	currentOperation := (*s)[index]
	(*s)[index-1].buffer = (*s)[index-1].Compute(currentOperation.buffer, (*s)[index-1].buffer)
	*s = (*s)[0:index]
}

func main() {
	inputFile := os.Args[1]
	dat, _ := ioutil.ReadFile(inputFile)

	expressions := strings.Split(strings.Replace(string(dat), " ", "", -1), "\n")
	totalSum := 0

	for _, exp := range expressions {
		operationStack := Stack([]Operation{Operation{"add", 0}})

		for _, opRune := range exp {
			switch opRune {
			case '(':
				operationStack.addOperation(`add`)
			case ')':
				operationStack.popOperation()
				for operationStack.Last().operation == `multiply` {
					operationStack.popOperation()
				}
			case '+':
				operationStack.Last().operation = `add`
			case '*':
				operationStack.Last().operation = `multiply`
				operationStack.addOperation(`add`)
			default:
				val, _ := strconv.Atoi(string(opRune))
				operationStack.Last().buffer = operationStack.Last().Compute(operationStack.Last().buffer, val)
			}
		}
		for len(operationStack) > 1 {
			operationStack.popOperation()
		}
		totalSum += operationStack[0].buffer
	}

	fmt.Println(totalSum)
}
