package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var DEBUG = true

func Debugf(msg string, subs ...any) {
	if DEBUG {
		log.Printf(msg, subs...)
	}
}

type Range64 struct {
	Start int64
	End   int64
}

func (r Range64) Contains(n int64) bool {
	return n >= r.Start && n <= r.End
}

func (r Range64) Overlaps(r2 Range64) bool {
	return r.Contains(r2.Start) || r.Contains(r2.End) || r2.Contains(r.Start) || r2.Contains(r.End)
}

func (r Range64) Count() int64 {
	return r.End - r.Start + 1
}

func (r Range64) String() string {
	return fmt.Sprintf("%d-%d", r.Start, r.End)
}

type Ranges []Range64

func (rs Ranges) Contains(n int64) bool {
	for _, r := range rs {
		if r.Contains(n) {
			return true
		}
	}
	return false
}

func (rs Ranges) Count() int64 {
	var total int64 = 0
	for _, r := range rs {
		total += r.Count()
	}
	return total
}

func (rs *Ranges) Reduce() {
	found := true
	for found {
		found = false
		for i, r1 := range *rs {
			var prune []int
			for j := i + 1; j < len(*rs); j++ {
				r2 := (*rs)[j]
				if r2.Overlaps(r1) {
					Debugf("Merging range %s into %s", r2.String(), r1.String())
					r1.Start = min(r1.Start, r2.Start)
					r1.End = max(r1.End, r2.End)
					Debugf("  Result: %s", r1.String())
					(*rs)[i] = r1
					prune = append([]int{j}, prune...)
				}
			}

			if len(prune) > 0 {
				for _, j := range prune {
					*rs = slices.Delete(*rs, j, j+1)
				}
				found = true
				// break out of the `for i, r1 := ...` loop since we've modified the
				// slice it's iterating over.
				break
			}
		}
	}
}

func (rs *Ranges) Add(newRange Range64) {
	for i, r := range *rs {
		if r.Overlaps(newRange) {
			Debugf("Merging range %s into %s", newRange.String(), r.String())
			r.Start = min(r.Start, newRange.Start)
			r.End = max(r.End, newRange.End)
			Debugf("  Result: %s", r.String())
			(*rs)[i] = r
			return
		}
	}
	Debugf("Could not merge %s into any range, appending", newRange.String())
	*rs = append(*rs, newRange)
}

func (rs Ranges) Dump() {
	for _, r := range rs {
		fmt.Printf("%s", r.String())
	}
}

func part1(data Ranges, ingredients []int64) int64 {
	var count int64 = 0
	for _, ingredient := range ingredients {
		if data.Contains(ingredient) {
			count++
		}
	}
	return count
}

func part2(data Ranges) int64 {
	return data.Count()
}

func parseLine(line string) Range64 {
	vals := strings.Split(line, "-")
	start, err := strconv.ParseInt(vals[0], 10, 64)
	if err != nil {
		log.Fatalf("Malformed input %s: %v", line, err)
	}
	end, err := strconv.ParseInt(vals[1], 10, 64)
	if err != nil {
		log.Fatalf("Malformed input %s: %v", line, err)
	}
	return Range64{Start: start, End: end}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var data Ranges
	var ingredients []int64
	inIngredients := false
	Debugf("===== PARSING INPUT =====")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if inIngredients {
			val, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				log.Fatalf("Malformed input %s: %v", line, err)
			}
			ingredients = append(ingredients, val)
		} else if len(line) == 0 {
			inIngredients = true
		} else {
			// data.Add(parseLine(line))
			data = append(data, parseLine(line))
		}
	}
	if DEBUG {
		Debugf("Parsed ranges:")
		data.Dump()
	}
	Debugf("===== REDUCING RANGES =====")
	data.Reduce()
	if DEBUG {
		Debugf("Reduced ranges:")
		data.Dump()
	}

	Debugf("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(data, ingredients)
	duration1 := time.Now().UnixMicro() - start1

	Debugf("===== PART 2 =====")
	start2 := time.Now().UnixMicro()
	part2Solution := part2(data)
	duration2 := time.Now().UnixMicro() - start2

	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("  in %dμs\n", duration1)
	fmt.Printf("Part 2: %d\n", part2Solution)
	fmt.Printf("  in %dμs\n", duration2)
}
