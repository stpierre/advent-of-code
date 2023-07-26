package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func sum(arr []int) int {
	sum := 0
	for _, valueInt := range arr {
		sum += valueInt
	}
	return sum
}

func main() {
	windowSizePtr := flag.Int("window", 1, "Size of the sliding window over which to consider depths")
	flag.Parse()

	depths := make([][]int, *windowSizePtr)
	count := 0
	readCount := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		current, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			log.Fatalf("Malformed input: %s", scanner.Text())
		}
		if readCount < *windowSizePtr {
			for i := 0; i <= readCount; i++ {
				depths[i] = append(depths[i], current)
			}
		} else {
			last := depths[0]
			for i := 0; i < *windowSizePtr-1; i++ {
				depths[i] = depths[i+1]
			}
			depths[*windowSizePtr-1] = make([]int, *windowSizePtr)

			for i := 0; i < *windowSizePtr; i++ {
				depths[i] = append(depths[i], current)
			}
			if sum(depths[0]) > sum(last) {
				count++
			}
		}
		readCount++
	}
	fmt.Println(count)
}
