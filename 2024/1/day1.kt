import kotlin.io.path.Path
import kotlin.io.path.readText
import kotlin.math.abs

/**
 * Reads lines from the given input txt file.
 */
fun readInput() = Path("input.txt").readText().trim().lines()

fun part1(lhs: MutableList<Int>, rhs: MutableList<Int>): Int {
    lhs.sort()
    rhs.sort()

    var sum = 0
    lhs.forEachIndexed { index, lhsValue ->
        sum += abs(lhsValue - rhs[index])
    }

    return sum
}

fun part2(lhs: List<Int>, rhs: List<Int>): Int {
    var sum = 0
    lhs.forEach { lhsValue ->
        sum += lhsValue * rhs.count { it == lhsValue }
    }

    return sum
}

fun main(args: Array<String>) {
    val lhs = mutableListOf<Int>()
    val rhs = mutableListOf<Int>()
    for (line in readInput()) {
        val values = line.split("\\s+".toRegex())
        lhs.add(values[0].toInt())
        rhs.add(values[1].toInt())
    }

    println("Part 1: ${part1(lhs, rhs)}")
    println("Part 2: ${part2(lhs, rhs)}")
}
