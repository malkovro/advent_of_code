package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Range struct {
	Low, High int
}

type FieldRule struct {
	RangeA, RangeB Range
}

type RuleSet []FieldRule

func (rs RuleSet) ValidRuleIndexes(i int) (compliantRuleIndex []int) {
	for ri, r := range rs {
		if r.IsValid(i) {
			compliantRuleIndex = append(compliantRuleIndex, ri)
		}
	}
	return
}

func (rs RuleSet) HasAnyValid(i int) bool {
	for _, r := range rs {
		if r.IsValid(i) {
			return true
		}
	}
	return false
}

func (r FieldRule) IsValid(i int) bool {
	return r.RangeA.Low <= i && r.RangeA.High >= i || r.RangeB.Low <= i && r.RangeB.High >= i
}

type Ticket []int

func (t Ticket) InvalidValues(rules RuleSet) (invalids []int) {
	for _, field := range t {
		if !rules.HasAnyValid(field) {
			invalids = append(invalids, field)
		}
	}
	return
}

func toInt(s string) int {
	sI, _ := strconv.Atoi(s)
	return sI
}

func processInput() (rules []FieldRule, tickets []Ticket, myTicket Ticket) {
	rulesRegex := regexp.MustCompile(`[^:]*: (?P<lowA>\d+)-(?P<highA>\d+) or (?P<lowB>\d+)-(?P<highB>\d+)$`)
	inputFile := os.Args[1]
	dat, _ := ioutil.ReadFile(inputFile)
	splittedInputs := strings.Split(string(dat), "\n\n")
	for _, line := range strings.Split(splittedInputs[0], "\n") {
		match := rulesRegex.FindStringSubmatch(line)
		rules = append(rules, FieldRule{Range{toInt(match[1]), toInt(match[2])}, Range{toInt(match[3]), toInt(match[4])}})
	}

	splittedTicketDefinition := strings.Split(splittedInputs[1], "\n")
	for _, iStr := range strings.Split(splittedTicketDefinition[1], ",") {
		i, _ := strconv.Atoi(iStr)
		myTicket = Ticket(append([]int(myTicket), i))
	}

	for i, line := range strings.Split(splittedInputs[2], "\n") {
		if i == 0 {
			continue
		}
		ticketDef := []int{}

		for _, iStr := range strings.Split(line, ",") {
			i, _ := strconv.Atoi(iStr)
			ticketDef = append(ticketDef, i)
		}
		tickets = append(tickets, Ticket(ticketDef))
	}

	return
}

func includes(list []int, item int) bool {
	for _, element := range list {
		if element == item {
			return true
		}
	}
	return false
}

func reduceComplianceMap(m *map[int][]int) {
	singletons := []int{}

	somethingChanged := true
	for somethingChanged {
		somethingChanged = false
		for k, compliantIndexes := range *m {
			if len(compliantIndexes) == 1 && !includes(singletons, compliantIndexes[0]) {
				somethingChanged = true
				singletons = append(singletons, compliantIndexes[0])
			} else if len(compliantIndexes) > 1 {
				newValues := []int{}
				for _, index := range compliantIndexes {
					if !includes(singletons, index) {
						newValues = append(newValues, index)
						somethingChanged = true
					}
				}
				(*m)[k] = newValues
			}
		}
	}
}

func getCompliancePerFieldMap(tickets []Ticket, ruleSet RuleSet) map[int][]int {
	compliancePerField := make(map[int][]int)
	for _, ticket := range tickets {
		for i, field := range ticket {
			ruleIndexes, exists := compliancePerField[i]
			compliantIndexesForField := []int{}
			if !exists {
				compliantIndexesForField = ruleSet.ValidRuleIndexes(field)
			} else {
				for _, r := range ruleIndexes {
					if ruleSet[r].IsValid(field) {
						compliantIndexesForField = append(compliantIndexesForField, r)
					}
				}
			}
			compliancePerField[i] = compliantIndexesForField
		}
	}

	reduceComplianceMap(&compliancePerField)
	return compliancePerField
}

func getIndexForField(i int, fieldMap map[int][]int) int {
	for k, v := range fieldMap {
		if v[0] == i {
			return k
		}
	}
	panic("Oops could  not find the field")
}

type Summable []int

func (list Summable) sum() (sum int) {
	for _, item := range list {
		sum += item
	}
	return
}

func main() {
	validTickets := []Ticket{}
	invalidValues := 0
	rules, tickets, myTicket := processInput()
	for _, ticket := range tickets {
		invalidValuesForTicket := ticket.InvalidValues(RuleSet(rules))
		invalidValues += Summable(invalidValuesForTicket).sum()
		if len(invalidValuesForTicket) == 0 {
			validTickets = append(validTickets, ticket)
		}
	}

	fmt.Println("Problem 1 => Sum of invalid values:", invalidValues)

	compliancePerField := getCompliancePerFieldMap(validTickets, RuleSet(rules))

	multiplied6FirstFields := 1
	for i, _ := range make([]int, 6) {
		index := getIndexForField(i, compliancePerField)
		multiplied6FirstFields *= myTicket[index]
	}

	fmt.Println("Problem 2 => The mulitplied 6 first field on my ticket give:", multiplied6FirstFields)
}
