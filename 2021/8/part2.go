package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

/* DIGIT  SEGMENTS  # SEGMENTS
 * 0      abcefg    6
 * 1      cf        2
 * 2      acdeg     5
 * 3      acdfg     5
 * 4      bcdf      4
 * 5      abdfg     5
 * 6      abdefg    6
 * 7      acf       3
 * 8      abcdefg   7
 * 9      abcdfg    6
 */

type Segment int

const (
	Top         Segment = 0
	TopLeft     Segment = 1
	TopRight    Segment = 2
	Middle      Segment = 3
	BottomLeft  Segment = 4
	BottomRight Segment = 5
	Bottom      Segment = 6
)

func findMissingSegment(digit1 string, digit2 string) rune {
	segments := make(map[rune]bool)
	for _, char := range digit1 {
		segments[char] = true
	}
	for _, char := range digit2 {
		segments[char] = false
	}
	for key, val := range segments {
		if val {
			return key
		}
	}
	return 0
}

func hasAllSegments(digit string, segments ...rune) bool {
	found := 0
	for _, candidate := range digit {
		for _, target := range segments {
			if candidate == target {
				found++
			}
		}
	}
	return found == len(segments)
}

func sortChars(chars string) string {
	charSlice := strings.Split(chars, "")
	sort.Strings(charSlice)
	return strings.Join(charSlice, "")
}

func getMapping(inputs []string) map[string]int {
	var wires [7]rune
	var digits [10]string

	// first, identify the easy digits: 1, 4, 7, and 8
	for _, input := range inputs {
		if len(input) == 2 {
			digits[1] = input
			fmt.Println(input, "is 1")
		} else if len(input) == 3 {
			digits[7] = input
			fmt.Println(input, "is 7")
		} else if len(input) == 4 {
			digits[4] = input
			fmt.Println(input, "is 4")
		} else if len(input) == 7 {
			digits[8] = input
			fmt.Println(input, "is 8")
		}
	}

	// identify the top wire: the segment of 7 that's not included in 1
	wires[Top] = findMissingSegment(digits[7], digits[1])
	fmt.Println("Wire", string(wires[Top]), "illuminates the top segment")

	// identify the 3: the only digit with 5 segments that includes both
	// segments from the digit 1
	for _, input := range inputs {
		if len(input) == 5 && hasAllSegments(input, []rune(digits[1])...) {
			digits[3] = input
			fmt.Println(input, "is 3")
			break
		}
	}

	// identify the top left wire: the segment of 4 that's not included in 3
	wires[TopLeft] = findMissingSegment(digits[4], digits[3])
	fmt.Println("Wire", string(wires[TopLeft]), "illuminates the top left segment")

	// identify the middle wire: the segment of 4 that's not included in
	// 1 and also is not the top left segment
	wires[Middle] = findMissingSegment(digits[4], string(append([]rune(digits[1]), wires[TopLeft])))
	fmt.Println("Wire", string(wires[Middle]), "illuminates the middle segment")

	/* identify the 2, 5, and 0:
	 * - 5 is the digit with 5 segments that includes the top left segment
	 * - 2 is the digit with 5 segments that is not 5 or 3
	 * - 0 is the digit with 6 segments that does not include the middle segment	 */
	for _, input := range inputs {
		if len(input) == 5 && input != digits[3] {
			if hasAllSegments(input, wires[TopLeft]) {
				digits[5] = input
				fmt.Println(input, "is 5")
			} else {
				digits[2] = input
				fmt.Println(input, "is 2")
			}
		} else if len(input) == 6 && !hasAllSegments(input, wires[Middle]) {
			digits[0] = input
			fmt.Println(input, "is 0")
		}
	}

	// identify the bottom left wire: the segment of 2 that's not included in 3
	wires[BottomLeft] = findMissingSegment(digits[2], digits[3])
	fmt.Println("Wire", string(wires[BottomLeft]), "illuminates the bottom left segment")

	/* identify the remaining numbers (6 and 9):
	 * - 6 is the digit with 6 segments that includes the top left segment
	 * - 9 is the digit with 6 segments that is not 6 or 0
	 */
	for _, input := range inputs {
		if len(input) == 6 && input != digits[0] {
			if hasAllSegments(input, wires[BottomLeft]) {
				digits[6] = input
				fmt.Println(input, "is 6")
			} else {
				digits[9] = input
				fmt.Println(input, "is 9")
			}
		}
	}

	// identify the bottom right wire: the segment of 3 that's not included in 2
	// wires[BottomRight] = findMissingSegment(digits[3], digits[2])
	// fmt.Println("Wire", string(wires[BottomRight]), "illuminates the bottom right segment")

	// identify the top right wire: the segment of 3 that's not included in 5
	// wires[TopLeft] = findMissingSegment(digits[3], digits[5])
	// fmt.Println("Wire", string(wires[TopRight]), "illuminates the top right segment")

	// identify the bottom wire: the only segment not yet identified
	// wires[Bottom] = findMissingSegment(digits[8], string(wires[:]))
	// fmt.Println("Wire", string(wires[Bottom]), "illuminates the bottom segment")

	retval := make(map[string]int)
	for digit, segments := range digits {
		retval[sortChars(segments)] = digit
	}
	return retval
}

func solve(inputs []string, outputs []string) int {
	mapping := getMapping(inputs)
	outputNum := 0
	for i, segments := range outputs {
		digit := mapping[sortChars(segments)]
		digitValue := digit * int(math.Pow(10, float64(len(outputs)-i-1)))
		fmt.Println("Converted", segments, "to", digit, "with value", digitValue)
		outputNum += digitValue
	}
	fmt.Println("Got output for", outputs, ":", outputNum)
	return outputNum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	total := 0
	for scanner.Scan() {
		data := strings.Split(strings.TrimSpace(scanner.Text()), "|")
		inputs := strings.Fields(strings.TrimSpace(data[0]))
		outputs := strings.Fields(strings.TrimSpace(data[1]))

		total += solve(inputs, outputs)
		fmt.Println("-------------------")
		fmt.Println("")
	}
	fmt.Println(total)
}
