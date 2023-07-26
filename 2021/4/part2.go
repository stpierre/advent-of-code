package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func indexOf(data []int, key int) int {
	for i, v := range data {
		if v == key {
			return i
		}
	}
	return -1
}

func findQuickestWin(board [5][5]int, numbers []int) int {
	boardMoves := len(numbers) + 1
	for row := 0; row < 5; row++ {
		rowMovesToWin := 0
		for col := 0; col < 5; col++ {
			rowMovesToWin = int(math.Max(
				float64(rowMovesToWin),
				float64(indexOf(numbers, board[row][col])),
			))
		}
		fmt.Println("  Row", row, "wins in", rowMovesToWin)
		if rowMovesToWin < boardMoves {
			boardMoves = rowMovesToWin
		}
	}

	for col := 0; col < 5; col++ {
		colMovesToWin := 0
		for row := 0; row < 5; row++ {
			colMovesToWin = int(math.Max(
				float64(colMovesToWin),
				float64(indexOf(numbers, board[row][col])),
			))
		}
		fmt.Println("  Column", col, "wins in", colMovesToWin)
		if colMovesToWin < boardMoves {
			boardMoves = colMovesToWin
		}
	}

	fmt.Println("Board wins in", boardMoves, "moves")

	return boardMoves
}

func sum(numbers [5]int) int {
	retval := 0
	for _, v := range numbers {
		retval += v
	}
	return retval
}

func findScore(board [5][5]int, numbers []int, numMoves int) int {
	var rowVals [5]int
	var allBoardNums []int
	for i, row := range board {
		rowVals[i] = sum(row)
		for _, v := range row {
			allBoardNums = append(allBoardNums, v)
		}
	}
	boardVal := sum(rowVals)

	for i := 0; i <= numMoves; i++ {
		if indexOf(allBoardNums, numbers[i]) != -1 {
			boardVal -= numbers[i]
		}
	}

	fmt.Println("Board value after", numMoves, "moves:", boardVal)
	boardScore := boardVal * numbers[numMoves]
	fmt.Println("Board score after", numMoves, "moves:", boardScore)
	return boardScore
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	var numbers []int
	for _, v := range strings.Split(strings.TrimSpace(scanner.Text()), ",") {
		num, _ := strconv.Atoi(v)
		numbers = append(numbers, num)
	}

	scanner.Scan() // read empty line after numbers

	var board [5][5]int
	maxMoves := 0
	bestScore := 0
	for rowNum := 0; scanner.Scan(); rowNum++ {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			boardMoves := findQuickestWin(board, numbers)
			if boardMoves > maxMoves {
				maxMoves = boardMoves
				bestScore = findScore(board, numbers, boardMoves)
			}

			// rowNum is incremented after the end of the loop, so we need
			// to set it to -1 so that it gets set to 0 for the next
			// iteration
			rowNum = -1
		} else {
			for i, v := range strings.Fields(line) {
				board[rowNum][i], _ = strconv.Atoi(v)
			}
		}
	}

	boardMoves := findQuickestWin(board, numbers)

	if boardMoves > maxMoves {
		maxMoves = boardMoves
		bestScore = findScore(board, numbers, boardMoves)
	}

	fmt.Println(bestScore)
}
