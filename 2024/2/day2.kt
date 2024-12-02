import kotlin.io.path.Path
import kotlin.io.path.readText
import kotlin.math.abs

fun debug(msg: String) {
    if (false) {
        println(msg)
    }
}

fun readInput() = Path("input.txt").readText().trim().lines()

fun isSafe(levels: List<Int>): Boolean {
    val increasing = levels[0] < levels[1]
    debug("$levels (increasing=$increasing):")
    levels.drop(1).forEachIndexed { idx, current ->
        // we are iterating over a copy of `levels` that has been
        // left-shifted by one -- the 0th item was removed and
        // everything else was shifted over -- so the 0th item in our
        // new iterable list is the 1st item in `levels`. when we
        // compare the 0th item in our iterable list and the 0th item
        // in `levels`, we're actually comparing the 1st item in
        // `levels` (`current`) with the 0th item in `levels`
        // (`previous`). are there less confusing ways to do this?
        // almost certainly.
        val previous = levels[idx]
        val delta = current - previous
        debug("  checking current=$current, previous=$previous, delta=$delta")
        if ((increasing && delta < 0) || (!increasing && delta > 0)) {
            debug("  not safe due to delta $delta")
            return false
        }
        if (abs(delta) < 1 || abs(delta) > 3) {
            debug("  not safe due to absolute delta ${abs(delta)}")
            return false
        }
    }

    debug("  verdict: SAFE!")
    return true
}

fun part1(reports: List<List<Int>>): Int {
    var safe = 0
    reports.forEach { levels ->
        if (isSafe(levels)) {
            safe++
        }
    }

    return safe
}

fun part2(reports: List<List<Int>>): Int {
    var safe = 0
    reports.forEach { levels ->
        if (isSafe(levels)) {
            safe++
        } else {
            val allIndices = (0..levels.size - 1).toSet()
            for (idx in levels.indices) {
	        val keep = allIndices.subtract(setOf(idx))
    	        if (isSafe(levels.slice(keep))) {
                    safe++
                    break
                }
            }
        }
    }

    return safe
}

fun main(args: Array<String>) {
    val reports = readInput().map { line ->
        line.split("\\s+".toRegex()).map {
            it.toInt()
        }
    }

    println("Part 1: ${part1(reports)}")
    println("Part 2: ${part2(reports)}")
}
