package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Heightmap struct {
	heights [][]int
	visited [][]bool
}

func (h *Heightmap) IsVisited(x int, y int) bool {
	return y < len(h.visited) && x < len(h.visited[y]) && h.visited[y][x]
}

func (h *Heightmap) Explore(x int, y int) int {
	// fmt.Println("  Exploring", x, ",", y)
	if h.heights[y][x] == 9 || h.IsVisited(x, y) {
		return 0
	}

	size := 1
	h.visited[y][x] = true

	if y-1 >= 0 {
		size += h.Explore(x, y-1)
	}

	if y+1 < len(h.heights) {
		size += h.Explore(x, y+1)
	}

	if x-1 >= 0 {
		size += h.Explore(x-1, y)
	}

	if x+1 < len(h.heights[y]) {
		size += h.Explore(x+1, y)
	}

	return size
}

func (h *Heightmap) InitializeVisited() {
	for y := 0; y < len(h.heights); y++ {
		h.visited = append(h.visited, []bool{})
		for x := 0; x < len(h.heights[y]); x++ {
			h.visited[y] = append(h.visited[y], false)
		}
	}
}

func main() {
	heightmap := Heightmap{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		charHeights := []rune(strings.TrimSpace(scanner.Text()))
		var lineHeights []int
		for _, v := range charHeights {
			height, _ := strconv.Atoi(string(v))
			lineHeights = append(lineHeights, height)
		}
		heightmap.heights = append(heightmap.heights, lineHeights)
	}
	heightmap.InitializeVisited()

	basins := []int{0, 0, 0}
	for y := 0; y < len(heightmap.heights); y++ {
		for x := 0; x < len(heightmap.heights[y]); x++ {
			fmt.Println("Starting exploration at", x, ",", y)
			basinSize := heightmap.Explore(x, y)
			for i, existing := range basins {
				if basinSize > existing {
					fmt.Println("Found new top 3 basin of size", basinSize)
					basins[i] = basinSize
					sort.Ints(basins)
					break
				}
			}
		}
	}

	fmt.Println("Found top 3 basins:", basins)
	fmt.Println(basins[0] * basins[1] * basins[2])
}
