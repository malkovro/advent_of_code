package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

type Rule struct {
	Id         string
	Match      string
	SubRuleIds [][]string
}

const maxRecursion = 50

var ruleRegex = regexp.MustCompile(`\d+`)

func (r Rule) ToPseudoRegex() string {
	if r.Match != "" {
		return r.Match
	}
	pseudoRegexes := []string{}
	for _, subRuleIds := range r.SubRuleIds {
		pseudoRegexes = append(pseudoRegexes, fmt.Sprintf("%v", strings.Join(subRuleIds, " ")))
	}
	return fmt.Sprintf("(%v)", strings.Join(pseudoRegexes, "|"))
}

func Resolve(ruleString *string, ruleMap map[string]Rule, visitedRules *map[string]int) (hasChanged bool) {
	*ruleString = ruleRegex.ReplaceAllStringFunc(*ruleString, func(s string) string {
		hasChanged = true
		(*visitedRules)[s]++
		switch s {
		case "8":
			if (*visitedRules)[s] > maxRecursion {
				return "(42)"
			}
		case "11":
			if (*visitedRules)[s] > maxRecursion {
				return "(42  31)"
			}
		}
		return ruleMap[s].ToPseudoRegex()
	})
	return
}

func main() {
	defer timeTrack(time.Now(), "Solving Time")
	inputFile := flag.String("f", "this is not optional!", "the input file")
	problem2Ptr := flag.Bool("problem2", false, "to solve the problem 2nd part")
	flag.Parse()

	problemNumber := 1
	if *problem2Ptr {
		problemNumber = 2
	}
	fmt.Printf("==> Solving for Problem %v (%v):\n", problemNumber, *inputFile)

	dat, _ := ioutil.ReadFile(*inputFile)

	lines := strings.Split(string(dat), "\n")
	count := 0

	ruleMap := make(map[string]Rule)
	ruleDefinitionDone := false
	visitedPaths := make(map[string]int)
	var ruleORegexp *regexp.Regexp

	for _, line := range lines {
		if line == "" {
			ruleDefinitionDone = true
			rule0 := ruleMap["0"].ToPseudoRegex()
			if *problem2Ptr {
				// Overwritte the rules for 8 & 11
				ruleMap["8"] = Rule{"8", "", [][]string{[]string{"42"}, []string{"42", "8"}}}
				ruleMap["11"] = Rule{"11", "", [][]string{[]string{"42", "31"}, []string{"42", "11", "31"}}}
			}

			hasChanged := true
			for hasChanged {
				hasChanged = Resolve(&rule0, ruleMap, &visitedPaths)
			}

			rule0 = fmt.Sprintf("^%v$", strings.Replace(rule0, " ", "", -1))
			ruleORegexp = regexp.MustCompile(rule0)
			continue
		}
		if ruleDefinitionDone {
			if ruleORegexp.MatchString(line) {
				count++
			}
			continue
		}

		rule := Rule{}

		for i, field := range strings.Fields(line) {
			if i == 0 {
				rule.Id = field[:len(field)-1]
				continue
			}

			if field == "|" {
				rule.SubRuleIds = append(rule.SubRuleIds, []string{})
				continue
			}
			if strings.HasPrefix(field, `"`) {
				rule.Match = strings.ReplaceAll(field, `"`, ``)
			} else {
				if len(rule.SubRuleIds) == 0 {
					rule.SubRuleIds = append(rule.SubRuleIds, []string{})
				}
				rule.SubRuleIds[len(rule.SubRuleIds)-1] = append(rule.SubRuleIds[len(rule.SubRuleIds)-1], field)
			}
		}
		ruleMap[rule.Id] = rule
	}

	fmt.Println("Number of valid messages", count)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
