package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func binToDec(digits []int) int {
	decimalNum := 0
	for i := len(digits) - 1; i >= 0; i-- {
		exponent := float64(len(digits) - i - 1)
		decimalNum += digits[i] * int(math.Pow(2, exponent))
	}
	return decimalNum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	firstLine := strings.TrimSpace(scanner.Text())
	digitSums := make([]int, len(firstLine))
	for i := 0; i < len(firstLine); i++ {
		bit, err := strconv.Atoi(string(firstLine[i]))
		if err != nil {
			log.Fatalf("Malformed input: %s", firstLine)
		}
		digitSums[i] = bit
	}

	readingCount := 1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		for i := 0; i < len(line); i++ {
			bit, err := strconv.Atoi(string(line[i]))
			if err != nil {
				log.Fatalf("Malformed input: %s", line)
			}
			digitSums[i] += bit
		}
		readingCount += 1
	}

	threshold := readingCount / 2
	fmt.Println("digit sums: ", digitSums)
	fmt.Println("threshold: ", threshold)

	gammaBits := make([]int, len(digitSums))
	for i := 0; i < len(digitSums); i++ {
		if digitSums[i] > threshold {
			gammaBits[i] = 1
		} else {
			gammaBits[i] = 0
		}
	}

	epsilonBits := make([]int, len(digitSums))
	for i := 0; i < len(gammaBits); i++ {
		epsilonBits[i] = 1 - gammaBits[i]
	}

	fmt.Println("gamma bits: ", gammaBits)
	fmt.Println("epsilon bits: ", epsilonBits)

	gamma := binToDec(gammaBits)
	epsilon := binToDec(epsilonBits)
	fmt.Println(gamma * epsilon)
}
