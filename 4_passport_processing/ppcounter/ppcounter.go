package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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

func homogenizePassportDefinition(passport string) string {
	return strings.Replace(passport, " ", "\n", -1)
}

func passportContainsAttribute(passport string, prefix string, attribute string, suffix string) bool {
	attributeRegexp := fmt.Sprintf("(?m)%v%v%v", prefix, attribute, suffix)
	matched, _ := regexp.Match(attributeRegexp, []byte(passport))

	return matched
}

func arePassportRequiredAttributesPresent(passport string) (bool, *string) {
	requiredAttributes := []string{"byr", `iyr`, `eyr`, `hgt`, `hcl`, `ecl`, `pid`}
	for _, attribute := range requiredAttributes {
		if !passportContainsAttribute(passport, `^`, attribute, `:`) {
			return false, &attribute
		}
	}
	return true, nil
}

func isPassportValid(passport string) (bool, *string) {
	attributeRegexps := []string{`byr:(19[2-8][0-9]|199[0-9]|200[0-2])`, `iyr:(201[0-9]|2020)`, `eyr:(20[12][0-9]|2030)`, `hgt:((1[5-8][0-9]|19[0-3])cm|(59|6[0-9]|7[0-6])in)`, `hcl:#([0-9]|[a-f]){6}`, `ecl:(amb|blu|brn|gry|grn|hzl|oth)`, `pid:\d{9}`}

	for _, attribute := range attributeRegexps {
		if !passportContainsAttribute(passport, `^`, attribute, `$`) {
			return false, &attribute
		}
	}
	return true, nil
}

func main() {
	inputFile := os.Args[1]
	batchPpt := readFile(inputFile)
	passports := strings.Split(batchPpt, "\n\n")
	passportWithReqAttributesCount := 0
	validishPassports := 0

	for _, passport := range passports {
		cleanedPassport := homogenizePassportDefinition(passport)
		if validity, err := arePassportRequiredAttributesPresent(cleanedPassport); validity {
			passportWithReqAttributesCount++
			if validity, err := isPassportValid(cleanedPassport); validity {
				validishPassports++
			} else {
				fmt.Println("Passport is invalid because missing: ", *err, "[", passport, "]")
			}
		} else {
			fmt.Println("Passport is invalid because missing: ", *err, "[", passport, "]")
		}
	}
	fmt.Println("Total number of  passports with all required attribute (except cid): ", passportWithReqAttributesCount)
	fmt.Println("Total number of  valid(ish) passports : ", validishPassports)
}
