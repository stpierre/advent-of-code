package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseLine(line string) (int, error) {
	direction := line[0]
	dist, err := strconv.Atoi(line[1:])
	if err != nil {
		return 0, err
	}
	if direction == 'L' {
		return dist * -1, nil
	}
	return dist, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pos := 50
	part1 := 0
	part2 := 0
	for scanner.Scan() {
		move, err := parseLine(strings.TrimSpace(scanner.Text()))
		if err != nil {
			log.Fatalf("Malformed input: %s", scanner.Text())
		}

		log.Printf("Processing move %s -> %d", scanner.Text(), move)
		rotations := int(math.Abs(float64(move / 100)))
		if rotations > 0 {
			log.Printf("  Move includes %d full rotations, adding to part 2", rotations)
			part2 += rotations
			log.Printf("  Move adjusted %d -> %d", move, move%100)
			move = move % 100
		}

		newPos := pos + move
		log.Printf("  Position %d -> %d (%d)", pos, newPos, newPos%100)
		if newPos >= 100 || (pos > 0 && newPos <= 0) {
			log.Printf("  Passed 0")
			part2++
		}
		pos = (newPos + 100) % 100
		if pos == 0 {
			log.Printf("  Landed on 0")
			part1++
		}
		log.Printf("  After move %s, position=%d, part1=%d, part2=%d",
			scanner.Text(), pos, part1, part2)
	}
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}
