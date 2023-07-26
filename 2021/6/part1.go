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
	daysPtr := flag.Int("days", 80, "Number of days to simulate")
	flag.Parse()

	var fish []int
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	sushi := strings.Split(strings.TrimSpace(scanner.Text()), ",")
	for _, strVal := range sushi {
		intVal, _ := strconv.Atoi(strVal)
		fish = append(fish, intVal)
	}

	for day := 1; day <= *daysPtr; day++ {
		var newFish []int
		for i := 0; i < len(fish); i++ {
			if fish[i] == 0 {
				fish[i] = 6
				newFish = append(newFish, 8)
			} else {
				fish[i]--
			}
		}
		fish = append(fish, newFish...)

		fmt.Println("After", day, "days:", fish)
	}

	fmt.Println(len(fish))
}
