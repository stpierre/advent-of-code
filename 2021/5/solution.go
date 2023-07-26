package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	skipDiagonalPtr := flag.Bool("skip-diagonals", false, "Skip diagonal vents")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	var grid [][]int
	for scanner.Scan() {
		coords := strings.Fields(strings.TrimSpace(scanner.Text()))
		startCoords := strings.Split(coords[0], ",")
		endCoords := strings.Split(coords[2], ",")
		if *skipDiagonalPtr && startCoords[0] != endCoords[0] && startCoords[1] != endCoords[1] {
			continue
		}

		fmt.Println("Found steam vent from", startCoords, "to", endCoords)

		startX, _ := strconv.Atoi(startCoords[0])
		endX, _ := strconv.Atoi(endCoords[0])
		var deltaX int
		if startX > endX {
			deltaX = -1
		} else if startX < endX {
			deltaX = 1
		} else {
			deltaX = 0
		}

		startY, _ := strconv.Atoi(startCoords[1])
		endY, _ := strconv.Atoi(endCoords[1])
		var deltaY int
		if startY > endY {
			deltaY = -1
		} else if startY < endY {
			deltaY = 1
		} else {
			deltaY = 0
		}

		// fmt.Println("  Iterating over", startX, "<= x <=", endX, "with delta =", deltaX)
		// fmt.Println("  Iterating over", startY, "<= y <=", endY, "with delta =", deltaY)

		x := startX
		y := startY
		for {
			for x >= len(grid) {
				grid = append(grid, []int{})
			}
			for y >= len(grid[x]) {
				grid[x] = append(grid[x], 0)
			}
			// fmt.Println("  Steam vent point at", x, ",", y)
			grid[x][y]++

			if startX != endX {
				x += deltaX
				if (deltaX > 0 && x > endX) || (deltaX < 0 && x < endX) {
					break
				}
			}
			if startY != endY {
				y += deltaY
				if (deltaY > 0 && y > endY) || (deltaY < 0 && y < endY) {
					break
				}
			}
		}

		// for x := 0; x < len(grid); x++ {
		// 	fmt.Println(grid[x])
		// }
	}

	overlaps := 0
	for x := 0; x < len(grid); x++ {
		for y := 0; y < len(grid[x]); y++ {
			if grid[x][y] > 1 {
				overlaps++
			}
		}
	}

	fmt.Println(overlaps)
}
