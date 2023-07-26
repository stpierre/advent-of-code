package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

func removeDuplicatePoints(points []Point) []Point {
	allKeys := make(map[Point]bool)
	var retval []Point
	for _, p := range points {
		if _, ok := allKeys[p]; !ok {
			allKeys[p] = true
			retval = append(retval, p)
		}
	}
	return retval
}

func printPage(points []Point) {
	xMax, yMax := 0, 0
	for _, p := range points {
		xMax = int(math.Max(float64(p.x), float64(xMax)))
		yMax = int(math.Max(float64(p.y), float64(yMax)))
	}

	var grid [][]bool
	for y := 0; y <= yMax; y++ {
		var row []bool
		for x := 0; x <= xMax; x++ {
			row = append(row, false)
		}
		grid = append(grid, row)
	}

	for _, p := range points {
		grid[p.y][p.x] = true
	}

	for y := 0; y <= yMax; y++ {
		for x := 0; x <= xMax; x++ {
			if grid[y][x] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func main() {
	foldRE := regexp.MustCompile(`fold along ([xy])=(\d+)`)
	scanner := bufio.NewScanner(os.Stdin)

	var points []Point
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, ",") {
			vals := strings.Split(line, ",")
			x, _ := strconv.Atoi(vals[0])
			y, _ := strconv.Atoi(vals[1])
			points = append(points, Point{x, y})
		} else if match := foldRE.FindStringSubmatch(line); match != nil {
			dimension := match[1]
			val, _ := strconv.Atoi(match[2])

			for i, p := range points {
				if dimension == "x" && p.x > val {
					p.x = val - (p.x - val)
					points[i] = p
				} else if dimension == "y" && p.y > val {
					p.y = val - (p.y - val)
					points[i] = p
				}
			}

			// fmt.Println("After fold:")
			// printPage(points)
			// fmt.Println()

			// uncomment this `break` for part1
			//break
		} else {
			// fmt.Println("All points read!")
			// printPage(points)
			// fmt.Println()
		}
	}

	fmt.Println("After fold:")
	printPage(points)
	fmt.Println()

	points = removeDuplicatePoints(points)
	fmt.Println(len(points))
}
