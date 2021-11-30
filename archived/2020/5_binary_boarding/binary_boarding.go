package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
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

func (boardingPass *BoardingPass) Row() (row int64) {
	rowBinaryStr := strings.Replace(boardingPass.Identifier, "F", "0", -1)
	rowBinaryStr = strings.Replace(rowBinaryStr, "B", "1", -1)
	row, _ = strconv.ParseInt(rowBinaryStr[0:7], 2, 64)
	return
}
func (boardingPass *BoardingPass) Column() (column int64) {
	columnBinaryStr := strings.Replace(boardingPass.Identifier, "L", "0", -1)
	columnBinaryStr = strings.Replace(columnBinaryStr, "R", "1", -1)
	column, _ = strconv.ParseInt(columnBinaryStr[7:10], 2, 64)
	return
}

type BoardingPass struct {
	Identifier string
}

func (boardingPass *BoardingPass) SeatId() int {
	return int(boardingPass.Column() + 8 * boardingPass.Row())
}

func arithmeticSum(u_n int, k int, r int) int {
	return	k * (u_n + (u_n - (k-1)* r)) / 2
}

func main() {
	inputFile := os.Args[1]
	batchBoardingPasses := readFile(inputFile)
	boardingPasses := strings.Split(batchBoardingPasses, "\n")
	summedSeatId := 0
	highestSeatId := 0

	for _, bp:= range boardingPasses {
		if len(bp) == 0 {
			fmt.Println("Highest Seat Id: ", highestSeatId)
			passengerCount := len(boardingPasses) - 1 + 1 // Removed the EOL but added myself!
			expectedSeatIdSum := arithmeticSum(highestSeatId, passengerCount, 1)
			fmt.Println("My Seat Id is: ", expectedSeatIdSum - summedSeatId)
			return
		}
		boardingPass := BoardingPass{bp}

		seatId := boardingPass.SeatId()
		summedSeatId = summedSeatId + seatId

		if seatId > highestSeatId {
			highestSeatId = seatId
		}
		// fmt.Println("BP ", boardingPass, " corresponds to seat ID", seatId, " on row: ", boardingPass.Row(), " and column: ", boardingPass.Column())
	}


}
