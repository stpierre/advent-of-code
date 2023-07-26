package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/* 0 abcefg  6
 * 1 cf      2
 * 2 acdeg   5
 * 3 acdfg   5
 * 4 bcdf    4
 * 5 abdfg   5
 * 6 abdefg  6
 * 7 acf     3
 * 8 abcdefg 7
 * 9 abcdfg  6
 */

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	count := 0
	for scanner.Scan() {
		outputs := strings.Fields(strings.TrimSpace(strings.Split(strings.TrimSpace(scanner.Text()), "|")[1]))

		for _, v := range outputs {
			if len(v) < 5 || len(v) == 7 {
				count++
			}
		}
	}
	fmt.Println(count)
}
