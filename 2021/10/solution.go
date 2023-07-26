package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(str string) {
	*s = append(*s, str)
}

func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

type Instruction struct {
	open              string
	close             string
	syntaxValue       int
	autocompleteValue int
}

func main() {
	instructions := map[string]Instruction{
		"(": Instruction{"(", ")", 3, 1},
		"[": Instruction{"[", "]", 57, 2},
		"{": Instruction{"{", "}", 1197, 3},
		"<": Instruction{"<", ">", 25137, 4},
	}

	var autocompleteScores []int
	var errors []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stack := Stack{}
		hasError := false
		line := strings.Split(strings.TrimSpace(scanner.Text()), "")
		for _, char := range line {
			_, ok := instructions[char]
			if ok { // opening instruction
				stack.Push(char)
			} else { // closing instruction
				opener, _ := stack.Pop()
				expected := instructions[opener]
				if expected.close != char {
					fmt.Println("Expected", expected.close, "but found", char, "instead")
					errors = append(errors, char)
					hasError = true
					break
				}
			}
		}

		if !hasError && len(stack) > 0 {
			fmt.Println("Found incomplete line", line)
			lineScore := 0
			for open, ok := stack.Pop(); ok; open, ok = stack.Pop() {
				instruction := instructions[open]
				lineScore = lineScore*5 + instruction.autocompleteValue
			}
			fmt.Println("  Autocomplete line value:", lineScore)
			autocompleteScores = append(autocompleteScores, lineScore)
		}
	}

	syntaxScore := 0
	for _, instruction := range instructions {
		count := 0
		for _, error := range errors {
			if error == instruction.close {
				count++
			}
		}

		fmt.Println("Found", count, "illegal", instruction.close, "instances for", count*instruction.syntaxValue, "points")
		syntaxScore += count * instruction.syntaxValue
	}

	fmt.Println("Syntax checker score:", syntaxScore)

	sort.Ints(autocompleteScores)
	fmt.Println("Autocomplete score:", autocompleteScores[len(autocompleteScores)/2])
}
