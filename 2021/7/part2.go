package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func fuelToMove(distance int) int {
	fuel := 0
	for i := 1; i <= distance; i++ {
		fuel += i
	}
	return fuel
}

func align(positions []int, alignAt int) int {
	fuel := 0
	for _, v := range positions {
		fuel += fuelToMove(int(math.Abs(float64(v - alignAt))))
	}
	return fuel
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	strPositions := strings.Split(strings.TrimSpace(scanner.Text()), ",")

	var positions []int
	maxPosition := 0
	minPosition := math.MaxInt32
	for _, strVal := range strPositions {
		intVal, _ := strconv.Atoi(strVal)
		positions = append(positions, intVal)
		maxPosition = int(math.Max(float64(maxPosition), float64(intVal)))
		minPosition = int(math.Min(float64(minPosition), float64(intVal)))
	}

	leastFuel := math.MaxInt32
	bestPosition := 0
	for candidate := minPosition; candidate <= maxPosition; candidate++ {
		fuel := align(positions, candidate)
		if fuel < leastFuel {
			leastFuel = fuel
			bestPosition = candidate
		}
	}

	fmt.Println("Align on", bestPosition, "consuming", leastFuel, "fuel")
}
