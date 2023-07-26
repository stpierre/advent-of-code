package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func filterBits(readings []string, index int, findMost bool) string {
	fmt.Println("round", index)
	ones := 0
	for _, reading := range readings {
		if reading[index] == '1' {
			ones++
		}
	}
	zeros := len(readings) - ones
	fmt.Println("  found", ones, "ones and", zeros, "zeros")

	var keep byte
	if (findMost && ones >= zeros) || (!findMost && zeros > ones) {
		keep = '1'
	} else {
		keep = '0'
	}
	fmt.Println("  keeping", string(keep))

	var newReadings []string
	for _, reading := range readings {
		if reading[index] == keep {
			newReadings = append(newReadings, reading)
		}
	}

	if len(newReadings) == 1 {
		fmt.Println("  only one option:", newReadings[0])
		return newReadings[0]
	}
	fmt.Println("  readings:", newReadings)
	return filterBits(newReadings, index+1, findMost)
}

func binToDec(digits string) int {
	decimalNum := 0
	for i := len(digits) - 1; i >= 0; i-- {
		exponent := float64(len(digits) - i - 1)
		digit, _ := strconv.Atoi(string(digits[i]))
		decimalNum += digit * int(math.Pow(2, exponent))
	}
	return decimalNum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var readings []string
	for scanner.Scan() {
		readings = append(readings, strings.TrimSpace(scanner.Text()))
	}

	oxygenGen := binToDec(filterBits(readings, 0, true))
	fmt.Println("oxygen generator rating:", oxygenGen)

	scrubber := binToDec(filterBits(readings, 0, false))
	fmt.Println("scrubber rating:", scrubber)

	fmt.Println("life support rating:", oxygenGen*scrubber)
}
