package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

var parentRegexp = regexp.MustCompile("^(?P<parent>(\\w*\\s)+)bags")

var cleanSpaceRegex = regexp.MustCompile("\\s")

func matchNamedRegexp(myRegexp *regexp.Regexp, str string) map[string]string {
	match := myRegexp.FindStringSubmatch(str)
	result := make(map[string]string)
	for i, name := range myRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[cleanSpaceRegex.ReplaceAllString(name, "")] = cleanSpaceRegex.ReplaceAllString(match[i], "")
		}
	}
	return result
}

type BagSet struct {
	Bag   string
	Count int
}

var cleanBagChildSet = regexp.MustCompile("(\\s|bag(s)?|\\.|no other)")
var childRegex = regexp.MustCompile("(?P<count>\\d+)(?P<child>\\w+)")

func childrenSet(str string) (bss []BagSet) {
	cleaned := cleanBagChildSet.ReplaceAllString(str, "")
	if cleaned == "" {
		return
	}

	for _, child := range strings.Split(cleaned, ",") {
		matches := matchNamedRegexp(childRegex, child)
		i, _ := strconv.Atoi(matches["count"])
		bss = append(bss, BagSet{matches["child"], i})
	}
	return
}

var childCleaningExp = regexp.MustCompile("(\\s|\\d+|bag(s)?|\\.|no other)")

func children(str string) []string {
	cleaned := childCleaningExp.ReplaceAllString(str, "")
	if cleaned == "" {
		return []string{}
	}
	return strings.Split(cleaned, ",")
}

func childBagsMap(rulesFile string) map[string][]BagSet {
	m := make(map[string][]BagSet)
	rules := strings.Split(rulesFile, "\n")
	for _, rule := range rules {
		if len(rule) == 0 {
			break
		}
		sections := strings.Split(rule, "contain")
		parent := matchNamedRegexp(parentRegexp, sections[0])["parent"]
		m[parent] = childrenSet(sections[1])
	}

	return m
}

func parentsBagMap(rulesFile string) map[string][]string {
	m := make(map[string][]string)
	rules := strings.Split(rulesFile, "\n")
	for _, rule := range rules {
		if len(rule) == 0 {
			break
		}
		sections := strings.Split(rule, "contain")
		parent := matchNamedRegexp(parentRegexp, sections[0])["parent"]
		for _, child := range children(sections[1]) {
			m[child] = append(m[child], parent)
		}
	}

	return m
}

type ParentList struct {
	Parents []string
}

func (parentList *ParentList) listParents(bag string, parentsMap map[string][]string) {
	directParents := parentsMap[bag]
	for _, parent := range directParents {
		parentList.Parents = append(parentList.Parents, parent)
		parentList.listParents(parent, parentsMap)
	}
}

func countChildBags(bag string, childBagsMap map[string][]BagSet) (count int) {
	for _, bs := range childBagsMap[bag] {
		count += bs.Count * (1 + countChildBags(bs.Bag, childBagsMap))
	}
	return
}

func main() {
	inputFile := os.Args[1]
	rulesFile := readFile(inputFile)

	parentsMap := parentsBagMap(rulesFile)
	childrenMap := childBagsMap(rulesFile)

	parentList := ParentList{}

	parentList.listParents("shinygold", parentsMap)

	fmt.Println("Shiny Gold can be contained in # many bags", len(unique(parentList.Parents)))
	fmt.Println("Shiny Gold must contain # many bags", countChildBags("shinygold", childrenMap))
}
