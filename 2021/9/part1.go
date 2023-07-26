package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	var heightmap [][]int
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		charHeights := []rune(strings.TrimSpace(scanner.Text()))
		var lineHeights []int
		for _, v := range charHeights {
			height, _ := strconv.Atoi(string(v))
			lineHeights = append(lineHeights, height)
		}
		heightmap = append(heightmap, lineHeights)
	}

	lowPointValues := 0
	for y := 0; y < len(heightmap); y++ {
		for x := 0; x < len(heightmap[y]); x++ {
			current := heightmap[y][x]

			above := math.MaxInt32
			if y-1 >= 0 {
				above = heightmap[y-1][x]
			}

			below := math.MaxInt32
			if y+1 < len(heightmap) {
				below = heightmap[y+1][x]
			}

			left := math.MaxInt32
			if x-1 >= 0 {
				left = heightmap[y][x-1]
			}

			right := math.MaxInt32
			if x+1 < len(heightmap[y]) {
				right = heightmap[y][x+1]
			}

			if current < above && current < below && current < left && current < right {
				fmt.Println("Found low point at", x, ",", y, ":", current)
				lowPointValues += current + 1
			}
		}
	}
	fmt.Println(lowPointValues)
}
