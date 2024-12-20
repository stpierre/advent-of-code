import kotlin.io.path.Path
import kotlin.io.path.readText

const val DEBUG = false

data class Equation(val total: Long, val inputs: List<Long>) {
    override fun toString(): String {
        return "${this.total}: ${this.inputs.joinToString(" ")}"
    }
}

fun readInput(): List<Equation> {
    return Path("input.txt").readText().trim().lines().map {
        val (total, numList) = it.split(": ")
        Equation(total.toLong(), numList.split(" ").map { it.toLong() })
    }
}


fun debug(msg: String) {
    if (DEBUG) {
        println(msg)
    }
}


fun getValues(vararg inputs: Long, allowConcat: Boolean = false): List<Long> {
    if (inputs.size == 1) {
        return listOf(inputs[0])
    } else {
        val results = mutableListOf<Long>()
        getValues(*inputs.slice(1..inputs.size - 1).toLongArray(), allowConcat = allowConcat).forEach {
            results.add(it * inputs[0])
            results.add(it + inputs[0])
            if (allowConcat) {
                results.add("${it}${inputs[0]}".toLong())
            }
        }
        return results
    }
}


fun solve(equations: List<Equation>, allowConcat: Boolean = false): Long {
    return equations.filter {
        if (DEBUG) {
            debug("${it} -> ${getValues(*it.inputs.reversed().toLongArray(), allowConcat = allowConcat)}")
        }
        getValues(*it.inputs.reversed().toLongArray(),
                  allowConcat = allowConcat).contains(it.total)
    }.sumOf {
        it.total
    }
}


fun main(args: Array<String>) {
    val equations = readInput()
    println("Part 1: ${solve(equations)}")
    println("Part 2: ${solve(equations, allowConcat = true)}")
}
