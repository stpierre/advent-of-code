package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var DEBUG = false

func Debugf(msg string, subs ...any) {
	if DEBUG {
		log.Printf(msg, subs...)
	}
}

type Range struct {
	Start int
	End   int
}

func (r Range) String() string {
	return fmt.Sprintf("%d-%d", r.Start, r.End)
}

func (r Range) FindNRepeats(n int) (repeats []int) {
	Debugf("Finding %d-repeats in %s", n, r.String())
	for i := r.Start; i <= r.End; i++ {
		val := strconv.Itoa(i)
		strLen := len(val)
		if strLen%n == 0 {
			matchLen := strLen / n
			found := true
		id:
			for firstIdx := range matchLen {
				for checkOffset := matchLen; checkOffset < strLen-firstIdx; checkOffset += matchLen {
					if val[firstIdx] != val[firstIdx+checkOffset] {
						found = false
						break id
					}
				}
			}

			if found {
				Debugf("  Found %d", i)
				repeats = append(repeats, i)
			}
		}
	}

	return repeats
}

func (r Range) FindAllRepeats() (repeats []int) {
	Debugf("Finding all repeats in %s", r.String())

	unique := map[int]bool{}
	for n := 2; n <= len(strconv.Itoa(r.End)); n++ {
		for _, i := range r.FindNRepeats(n) {
			unique[i] = true
		}
	}

	for k := range unique {
		repeats = append(repeats, k)
	}
	return repeats
}

func parseRange(rangeStr string) (Range, error) {
	vals := strings.Split(rangeStr, "-")
	start, err := strconv.Atoi(vals[0])
	if err != nil {
		return Range{}, err
	}
	end, err := strconv.Atoi(vals[1])
	if err != nil {
		return Range{}, err
	}
	return Range{Start: start, End: end}, nil
}

func part1(ranges []Range) int {
	sum := 0
	for _, idRange := range ranges {
		for _, repeat := range idRange.FindNRepeats(2) {
			sum += repeat
		}
	}
	return sum
}

func part2(ranges []Range) int {
	sum := 0
	for _, idRange := range ranges {
		for _, repeat := range idRange.FindAllRepeats() {
			sum += repeat
		}
	}
	return sum
}

func main() {
	var ranges []Range
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		for rangeStr := range strings.SplitSeq(strings.TrimSpace(scanner.Text()), ",") {
			idRange, err := parseRange(rangeStr)
			if err != nil {
				log.Fatalf("Malformed input: %v: %s", err, rangeStr)
			}
			ranges = append(ranges, idRange)
		}
	} else {
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading stdin: %v", err)
		} else {
			log.Fatalf("No input")
		}
	}

	fmt.Println("===== PART 1 =====")
	part1Solution := part1(ranges)
	fmt.Println("===== PART 2 =====")
	part2Solution := part2(ranges)
	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("Part 2: %d\n", part2Solution)
}
