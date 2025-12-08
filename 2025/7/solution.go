package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var DEBUG = false

func Debugf(msg string, subs ...any) {
	if DEBUG {
		log.Printf(msg, subs...)
	}
}

func part1(grid [][]bool, start int) int {
	splits := 0
	beams := make([]bool, len(grid[0]))
	beams[start] = true
	Debugf("At start, beam at %d", start)
	for rowNum, row := range grid {
		Debugf("Processing row %d", rowNum)
		newBeams := make([]bool, len(row))
		for col, hasBeam := range beams {
			if !hasBeam {
				continue
			}

			if row[col] {
				Debugf("  Splitting beam at %d", col)
				if col-1 >= 0 {
					newBeams[col-1] = true
				}
				if col+1 < len(row) {
					newBeams[col+1] = true
				}
				splits++
			} else {
				newBeams[col] = true
			}
		}
		if DEBUG {
			var beamsAt []int
			for col, hasBeam := range newBeams {
				if hasBeam {
					beamsAt = append(beamsAt, col)
				}
			}
			Debugf("  Beams at %v", beamsAt)
		}
		beams = newBeams
	}
	return splits
}

type Point struct {
	Row int
	Col int
}

var countPathsCache = make(map[Point]int)

func countPaths(grid *[][]bool, start int, startRow int) int {
	p := Point{Row: startRow, Col: start}
	if val, ok := countPathsCache[p]; ok {
		return val
	}

	Debugf("Counting paths from row %d", startRow)
	if startRow >= len(*grid) {
		return 1
	}

	paths := 0
	row := (*grid)[startRow]
	if row[start] {
		Debugf("  Splitting time at %d", start)
		if start-1 >= 0 {
			paths += countPaths(grid, start-1, startRow+1)
		}
		if start+1 < len(row) {
			paths += countPaths(grid, start+1, startRow+1)
		}
	} else {
		Debugf("  Continuing at %d", start)
		paths += countPaths(grid, start, startRow+1)
	}

	Debugf("For row %d, got %d paths", startRow, paths)
	countPathsCache[p] = paths
	return paths
}

func part2(grid *[][]bool, start int) int {
	return countPaths(grid, start, 0)
}

func parseLine(line string) (row []bool) {
	for _, char := range line {
		row = append(row, char == '^')
	}
	return row
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	start := strings.Index(scanner.Text(), "S")

	var grid [][]bool
	for scanner.Scan() {
		grid = append(grid, parseLine(strings.TrimSpace(scanner.Text())))
	}

	Debugf("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(grid, start)
	duration1 := time.Now().UnixMicro() - start1

	Debugf("===== PART 2 =====")
	start2 := time.Now().UnixMicro()
	part2Solution := part2(&grid, start)
	duration2 := time.Now().UnixMicro() - start2

	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("  in %dμs\n", duration1)
	fmt.Printf("Part 2: %d\n", part2Solution)
	fmt.Printf("  in %dμs\n", duration2)
}
