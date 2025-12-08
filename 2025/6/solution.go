package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var DEBUG = true

func Debugf(msg string, subs ...any) {
	if DEBUG {
		log.Printf(msg, subs...)
	}
}

func performOp(op rune, vals []int64) (result int64) {
	if op == '+' {
		for _, i := range vals {
			result += i
		}
	} else { // *
		result = 1
		for _, i := range vals {
			result *= i
		}
	}
	return result
}

func part1(numbers [][]int64, operations []rune) (result int64) {
	for idx, op := range operations {
		var column []int64
		for _, row := range numbers {
			column = append(column, row[idx])
		}
		val := performOp(op, column)
		Debugf("Got %d from %c on column %d", val, op, idx)
		result += val
	}
	return result
}

func part2(lines []string, operations []rune) (result int64) {
	var fieldLines [][]string
	for _, line := range lines {
		fieldLines = append(fieldLines, strings.Fields(line))
	}

	offset := 0
	for col, op := range operations {
		colWidth := 0
		for _, row := range fieldLines {
			colWidth = max(colWidth, len(row[col]))
		}
		Debugf("Column %d is %d characters wide", col, colWidth)

		numbers := make([]int64, colWidth)
		for _, line := range lines {
			for i := range colWidth {
				// the last column doesn't have trailing spaces to pad it out to the full length
				idx := min(i+offset, len(line)-1)
				if raw := line[idx]; raw != ' ' {
					val, err := strconv.ParseInt(string(raw), 10, 64)
					if err != nil {
						log.Fatalf("Malformed input %s (%s): %v", line, raw, err)
					}
					Debugf("Column %d number %d digit = %d", col, i, val)
					numbers[i] = numbers[i]*10 + val
				}
			}
		}
		Debugf("For column %d got numbers %v", col, numbers)

		val := performOp(op, numbers)
		Debugf("Got %d from %c on column %d", val, op, col)
		result += val

		offset += colWidth + 1
	}
	return result
}

func parseOpLine(line string) (result []rune) {
	for _, raw := range strings.Fields(line) {
		result = append(result, []rune(raw)[0])
	}
	return result
}

func parseDigitLinePart1(line string) (result []int64) {
	for _, raw := range strings.Fields(line) {
		val, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			log.Fatalf("Malformed input %s (%s): %v", line, raw, err)
		}
		result = append(result, val)
	}
	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	var numbers [][]int64
	var operations []rune
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line[0] == '+' || line[0] == '*' {
			operations = parseOpLine(line)
		} else {
			numbers = append(numbers, parseDigitLinePart1(line))
			lines = append(lines, strings.TrimRight(scanner.Text(), "\n"))
		}
	}

	Debugf("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(numbers, operations)
	duration1 := time.Now().UnixMicro() - start1

	Debugf("===== PART 2 =====")
	start2 := time.Now().UnixMicro()
	part2Solution := part2(lines, operations)
	duration2 := time.Now().UnixMicro() - start2

	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("  in %dμs\n", duration1)
	fmt.Printf("Part 2: %d\n", part2Solution)
	fmt.Printf("  in %dμs\n", duration2)
}
