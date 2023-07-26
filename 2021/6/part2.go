package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const Threshold = 80

func reduce(days int) []int {
	if days < Threshold {
		return []int{days}
	}
	return append(reduce(days-7), reduce(days-9)...)
}

func main() {
	daysPtr := flag.Int("days", 256, "Number of days to simulate")
	flag.Parse()

	days := reduce(*daysPtr)
	fmt.Println("Sum of days:", days)

	var fish []int
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	sushi := strings.Split(strings.TrimSpace(scanner.Text()), ",")
	for _, strVal := range sushi {
		intVal, _ := strconv.Atoi(strVal)
		fish = append(fish, intVal)
	}

	var counts [Threshold + 1]int
	counts[0] = len(fish)
	for day := 1; day <= Threshold; day++ {
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

		fmt.Println("After", day, "days:", len(fish))
		counts[day] = len(fish)
	}

	total := 0
	for _, v := range days {
		total += counts[v]
	}
	fmt.Println(total)
}
