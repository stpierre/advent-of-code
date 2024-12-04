import kotlin.io.path.Path
import kotlin.io.path.readText
import kotlin.math.abs

fun debug(msg: String) {
    if (true) {
        println(msg)
    }
}

// read as a single string, not an array of lines
fun readInput() = Path("input.txt").readText().trim()

fun part1(input: String): Int {
    val mulRE = Regex("mul\\((\\d+),(\\d+)\\)",
                      setOf(RegexOption.MULTILINE, RegexOption.DOT_MATCHES_ALL))
    var total = 0
    mulRE.findAll(input).forEach {
        total += it.groups.get(1)!!.value.toInt() * it.groups.get(2)!!.value.toInt()
    }
    return total
}

fun part2(input: String): Int {
    val instructionRE = Regex("do\\(\\)|don't\\(\\)|mul\\((\\d+),(\\d+)\\)",
                      setOf(RegexOption.MULTILINE, RegexOption.DOT_MATCHES_ALL))
    var total = 0
    var enabled = true
    instructionRE.findAll(input).forEach {
        val instruction = it.groups.get(0)!!.value
        if (instruction == "do()") {
            enabled = true
        } else if (instruction == "don't()") {
            enabled = false
        } else if (enabled) {
            total += it.groups.get(1)!!.value.toInt() * it.groups.get(2)!!.value.toInt()
        }
    }
    return total
}

fun main(args: Array<String>) {
    val input = readInput()
    println("Part 1: ${part1(input)}")
    println("Part 2: ${part2(input)}")
}
