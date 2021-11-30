package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type WriteInstruction struct {
	Address, Value int
}

type DockingMask string

func (m DockingMask) AsClearingByteMask() Bits {
	m = DockingMask(strings.Replace(string(m), "1", "X", -1))
	m = DockingMask(strings.Replace(string(m), "0", "1", -1))
	m = DockingMask(strings.Replace(string(m), "X", "0", -1))

	v, err := strconv.ParseUint(string(m), 2, 64)
	check(err)
	return Bits(v)
}

func (m DockingMask) AsSettingByteMask() Bits {
	m = DockingMask(strings.Replace(string(m), "X", "0", -1))
	v, err := strconv.ParseUint(string(m), 2, 64)

	check(err)
	return Bits(v)
}

func PowInt(a, n int) int { return int(math.Pow(float64(a), float64(n))) }

func (m DockingMask) LoopOverFloating() []DockingMask {
	floatingCount := strings.Count(string(m), "X")
	numberOfCases := PowInt(2, floatingCount)
	masks := make([]DockingMask, numberOfCases)
	formattingString := fmt.Sprintf("%%0%ds", floatingCount)
	baseMask := strings.Replace(string(m), "0", "%", -1)
	for i, _ := range make([]int, numberOfCases) {
		mask := baseMask
		floatingCaseBinary := fmt.Sprintf(formattingString, strconv.FormatInt(int64(i), 2))

		for index, _ := range make([]int, floatingCount) {
			mask = strings.Replace(mask, "X", string(floatingCaseBinary[index]), 1)
		}
		mask = strings.Replace(mask, "%", "X", -1)
		masks[i] = DockingMask(mask)
	}
	return masks
}

func extractWriteInstruction(str string) WriteInstruction {
	regex := regexp.MustCompile(`^mem\[(?P<address>\d+)\]\s=\s(?P<value>.*)$`)
	match := regex.FindStringSubmatch(str)
	add, _ := strconv.Atoi(match[1])
	val, _ := strconv.Atoi(match[2])
	return WriteInstruction{add, val}
}

func extractMask(str string) (DockingMask, error) {
	maskRegex := regexp.MustCompile(`^mask\s=\s(?P<mask>.*)$`)
	match := maskRegex.FindStringSubmatch(str)

	if len(match) == 0 {
		return DockingMask(""), errors.New("No mask in here!")
	}

	return DockingMask(match[1]), nil
}

type DockingMemory map[int]int

func (dm DockingMemory) sum() (sum int) {
	for _, v := range dm {
		sum += v
	}
	return
}

type Bits uint64

func Set(b, flag Bits) Bits   { return b | flag }
func Clear(b, flag Bits) Bits { return b &^ flag }

func main() {
	inputFile := os.Args[1]
	input := readFile(inputFile)
	inputData := strings.Split(input, "\n")

	var mask DockingMask
	var maskExtractionErr error
	var clearingByteMask Bits
	var settingByteMask Bits
	var floatingMasks []DockingMask
	memory1 := make(DockingMemory)
	memory2 := make(DockingMemory)

	for _, instruction := range inputData {
		mask, maskExtractionErr = extractMask(instruction)

		if maskExtractionErr == nil {
			floatingMasks = mask.LoopOverFloating()
			clearingByteMask = mask.AsClearingByteMask()
			settingByteMask = mask.AsSettingByteMask()
			continue
		}

		wi := extractWriteInstruction(instruction)
		value := Set(Bits(wi.Value), settingByteMask)
		value = Clear(value, clearingByteMask)
		memory1[wi.Address] = int(value)

		for _, floatingMask := range floatingMasks {
			address := Set(Bits(wi.Address), floatingMask.AsSettingByteMask())
			address = Clear(address, floatingMask.AsClearingByteMask())
			memory2[int(address)] = wi.Value
		}
	}

	fmt.Println("Problem 1: The total sum of in-memory values after value-masking is:", memory1.sum())
	fmt.Println("Problem 2: The total sum of in-memory values after address-masking is:", memory2.sum())
}
