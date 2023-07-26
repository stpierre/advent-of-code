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
	var octopodes [][]int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), "")
		var octopusRow []int
		for _, strVal := range line {
			intval, _ := strconv.Atoi(strVal)
			octopusRow = append(octopusRow, intval)
		}
		octopodes = append(octopodes, octopusRow)
	}

	totalOctopodes := len(octopodes) * len(octopodes[0])
	var flashes int
	var step int
	for step = 1; flashes < totalOctopodes; step++ {
		flashes = 0
		if step%100 == 0 {
			fmt.Println("Step", step)
		}

		var flashed [][]bool
		for y := 0; y < len(octopodes); y++ {
			var flashedRow []bool
			for x := 0; x < len(octopodes[y]); x++ {
				octopodes[y][x]++
				flashedRow = append(flashedRow, false)
			}

			flashed = append(flashed, flashedRow)
		}

		for true {
			found := false

			for y := 0; y < len(octopodes); y++ {
				for x := 0; x < len(octopodes[y]); x++ {
					if octopodes[y][x] > 9 && !flashed[y][x] {
						flashed[y][x] = true
						found = true
						flashes++
						yMin := int(math.Max(0.0, float64(y-1)))
						yMax := int(math.Min(float64(len(octopodes)-1), float64(y+1)))
						xMin := int(math.Max(0.0, float64(x-1)))
						xMax := int(math.Min(float64(len(octopodes[y])-1), float64(x+1)))
						for updateY := yMin; updateY <= yMax; updateY++ {
							for updateX := xMin; updateX <= xMax; updateX++ {
								if updateX != x || updateY != y {
									octopodes[updateY][updateX]++
								}
							}
						}
					}
				}
			}

			if !found {
				break
			}
		}

		for y := 0; y < len(octopodes); y++ {
			for x := 0; x < len(octopodes[y]); x++ {
				if flashed[y][x] {
					octopodes[y][x] = 0
				}
			}
		}
	}

	// subtract 1 because of the postincrement on the loop
	fmt.Println(step - 1)
}
