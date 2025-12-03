package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

var DEBUG = true

func Debugf(msg string, subs ...any) {
	if DEBUG {
		log.Printf(msg, subs...)
	}
}

type Bank []int

func (b Bank) String() string {
	result := make([]string, len(b))
	for i, v := range b {
		result[i] = strconv.Itoa(v)
	}
	return strings.Join(result, "")
}

func (b Bank) MaxJoltage(digits int) int64 {
	start := 0
	end := len(b) - (digits - 1)
	var result int64
	Debugf("Finding max joltage in %s", b.String())
	for i := range digits {
		val := slices.Max(b[start:end])
		Debugf("  Got %dth digit %d from start=%d, end=%d, range=%v",
			i, val, start, end, b[start:end])
		start = start + slices.Index(b[start:end], val) + 1
		end = len(b) - (digits - (i + 2))
		Debugf("  New start=%d, end=%d", start, end)
		result += int64(val) * int64(math.Pow(10, float64(digits-i-1)))
		Debugf("  Running total %d", result)
	}
	Debugf("  Got max joltage from %s: %d", b.String(), result)
	return result
}

func parseLine(line string) Bank {
	var result Bank
	for _, char := range line {
		result = append(result, int(char-'0'))
	}
	return result
}

func solve(banks []Bank, digits int) int64 {
	var totalJoltage int64
	for _, bank := range banks {
		totalJoltage += bank.MaxJoltage(digits)
	}
	return totalJoltage
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var banks []Bank
	for scanner.Scan() {
		banks = append(banks, parseLine(strings.TrimSpace(scanner.Text())))
	}

	fmt.Println("===== PART 1 =====")
	part1Solution := solve(banks, 2)
	fmt.Println("===== PART 2 =====")
	part2Solution := solve(banks, 12)
	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("Part 2: %d\n", part2Solution)
}
