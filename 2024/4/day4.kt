import kotlin.io.path.Path
import kotlin.io.path.readText

fun debug(msg: String) {
    if (false) {
        println(msg)
    }
}

fun readInput() = Path("test_input.txt").readText().trim().lines()

fun String.findAllIndices(match: String): List<Int> {
    var startIndex = 0
    val results = mutableListOf<Int>()
    while (true) {
        val index = this.indexOf(match, startIndex)
        if (index == -1) break

        results.add(index)
        startIndex = index + match.length
    }
    return results
}

fun getStep(start: Int, end: Int): Int {
    if (start > end) {
        return -1
    } else if (end > start) {
        return 1
    } else {
        return 0
    }
}

fun isXmas(input: List<String>, range: Range2D): Boolean {
    val goal = "XMAS"
    if (range.bounds.any { it < 0 } ||
            range.rowBounds.any { it >= input.size } ||
            range.colBounds.any { it >= input[0].length }) {
        return false
    }

    val found = (0..goal.length - 1).map { offset ->
        val row = range.startRow + offset * range.rowStep
        val col = range.startCol + offset * range.colStep
        input[row][col]
    }.joinToString("")
    return found == goal
}

class Range2D(val startRow: Int, val endRow: Int,
              val startCol: Int, val endCol: Int) {
    val rowStep = getStep(startRow, endRow)
    val colStep = getStep(startCol, endCol)

    val bounds = listOf(startRow, endRow, startCol, endCol)
    val rowBounds = listOf(startRow, endRow)
    val colBounds = listOf(startCol, endCol)
}

fun part1(input: List<String>): Long {
    return input.indices.sumOf { rowNum ->
        val row = input[rowNum]
        row.findAllIndices("X").sumOf { colNum ->
            listOf(
                // horizontal, forward
                Range2D(rowNum, rowNum, colNum, colNum + 3),
                // horizontal, backward
                Range2D(rowNum, rowNum, colNum, colNum - 3),
                // vertical, downward
                Range2D(rowNum, rowNum + 3, colNum, colNum),
                // vertical, upward
                Range2D(rowNum, rowNum - 3, colNum, colNum),
                // diagonal, up and to the right
                Range2D(rowNum, rowNum - 3, colNum, colNum + 3),
                // diagonal, down and to the right
                Range2D(rowNum, rowNum + 3, colNum, colNum + 3),
                // diagonal, down and to the left
                Range2D(rowNum, rowNum + 3, colNum, colNum - 3),
                // diagonal, up and to the left
                Range2D(rowNum, rowNum - 3, colNum, colNum - 3),
            ).sumOf {
                if (isXmas(input, it)) 1L else 0L
            }
        }
    }
}

fun part2(input: List<String>): Long {
    return input.indices.sumOf { rowNum ->
        val row = input[rowNum]
        row.findAllIndices("A").sumOf { colNum ->
            if (rowNum > 0 && colNum > 0 &&
                    rowNum + 1 < input.size && colNum + 1 < row.length) {
                val topLeft = input[rowNum - 1][colNum - 1]
                val topRight = input[rowNum - 1][colNum + 1]
                val bottomLeft = input[rowNum + 1][colNum - 1]
                val bottomRight = input[rowNum + 1][colNum + 1]
                if (((topLeft == 'M' && bottomRight == 'S') ||
                         (topLeft == 'S' && bottomRight == 'M')) &&
                        ((topRight == 'M' && bottomLeft == 'S') ||
                             (topRight == 'S' && bottomLeft == 'M'))) {
                    1L
                } else {
                    0L
                }
            } else {
                0L
            }
        }
    }
}

fun main(args: Array<String>) {
    val input = readInput()
    println("Part 1: ${part1(input)}")
    println("Part 2: ${part2(input)}")
}
