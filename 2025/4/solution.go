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

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d, %d", p.X, p.Y)
}

var adjustments = [8]Point{
	{X: -1, Y: -1},
	{X: -1, Y: 0},
	{X: -1, Y: 1},
	{X: 0, Y: -1},
	{X: 0, Y: 1},
	{X: 1, Y: -1},
	{X: 1, Y: 0},
	{X: 1, Y: 1},
}

type Grid [][]bool

func (g Grid) findAccessible() []Point {
	var points []Point
	threshold := 4
	for y, row := range g {
		for x, point := range row {
			if point {
				Debugf("Counting adjacent rolls at %d, %d", x, y)
				adjacent := 0
				for _, adj := range adjustments {
					newX, newY := x-adj.X, y-adj.Y
					if newX >= 0 && newX < len(row) &&
						newY >= 0 && newY < len(g) &&
						g[newY][newX] {
						Debugf("  Found at %d, %d", newX, newY)
						adjacent++
					}

					if adjacent >= threshold {
						Debugf("  Already found %d rolls, bailing out", adjacent)
						break
					}
				}

				Debugf("  Found %d rolls", adjacent)
				if adjacent < threshold {
					points = append(points, Point{X: x, Y: y})
				}
			}
		}
	}
	return points
}

func (g Grid) countAccessible() int {
	return len(g.findAccessible())
}

func part1(data Grid) int {
	return data.countAccessible()
}

func part2(data Grid) int {
	removed := 0
	points := data.findAccessible()
	for len(points) > 0 {
		for _, point := range points {
			Debugf("Removing roll at %s", point.String())
			removed++
			data[point.Y][point.X] = false
		}
		points = data.findAccessible()
	}
	return removed
}

func parseLine(line string) []bool {
	var result []bool
	for _, char := range line {
		result = append(result, char == '@')
	}
	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var data Grid
	for scanner.Scan() {
		data = append(data, parseLine(strings.TrimSpace(scanner.Text())))
	}

	Debugf("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(data)
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
