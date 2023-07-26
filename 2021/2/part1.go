package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	position := 0
	depth := 0
	for scanner.Scan() {
		words := strings.Fields(strings.TrimSpace(scanner.Text()))
		direction := words[0]
		distance, err := strconv.Atoi(words[1])
		if err != nil {
			log.Fatalf("Malformed input: %s", scanner.Text())
		}
		if direction == "forward" {
			position += distance
		} else if direction == "down" {
			depth += distance
		} else if direction == "up" {
			depth -= distance
		}
		fmt.Println(scanner.Text())
		fmt.Printf("depth=%d position=%d\n", depth, position)
	}
	fmt.Println(depth * position)
}
