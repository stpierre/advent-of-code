import kotlin.io.path.Path
import kotlin.io.path.readText

const val DEBUG = false


fun debug(msg: String) {
    if (DEBUG) {
        println(msg)
    }
}


fun blink(num: Long): Pair<Long, Long?> {
    val numAsString = num.toString()
    val digits = numAsString.length
    if (num == 0L) {
        return Pair(1L, null)
    } else if (digits % 2 == 0) {
        return Pair(numAsString.substring(0, digits / 2).toLong(),
                    numAsString.substring(digits / 2, digits).toLong())
    } else {
        return Pair(num * 2024, null)
    }
}


fun readInput(): List<Long> {
    return Path("input.txt").readText().trim().lines()[0].split("\\s+".toRegex()).map {
        it.toLong()
    }
}


fun solveNaive(stones: MutableList<Long>, count: Int = 25): Int {
    repeat(count) {
        var idx = 0
        while (idx < stones.size) {
            val (newVal, newStone) = blink(stones[idx])
            stones[idx] = newVal
            if (newStone != null) {
                stones.add(idx + 1, newStone!!)
                idx++
            }
            idx++
        }
    }
    return stones.size
}


fun solveUnordered(input: List<Long>, count: Int = 75): Long {
    val stones = mutableMapOf<Long, Long>()
    input.forEach {
        stones[it] = 1
    }
    repeat(count) {
        debug("Iteration $it")
        stones.toMap().forEach { stone, count ->
            val (newVal, newStone) = blink(stone)
            debug("  Replacing $count $stone stones with $newVal")
            if (stones[stone] == count) {
                stones.remove(stone)
                debug("  No $stone stones remain")
            } else {
                stones[stone] = stones[stone]!! - count
                debug("  ${stones[stone]} $stone stones remain")
            }
            stones.merge(newVal, count) { oldCount, newCount ->
                oldCount + newCount
            }
            debug("  ${stones[newVal]} total $newVal stones")
            if (newStone != null) {
                debug("  Adding $count $newStone stones")
                stones.merge(newStone, count) { oldCount, newCount ->
                    oldCount + newCount
                }
                debug("  ${stones[newStone]} total $newStone stones")
            }
        }
        debug("After iteration $it, stones=$stones (${stones.values.sumOf { it }})")
    }
    return stones.values.sumOf { it }
}


fun main(args: Array<String>) {
    val input = readInput()
    println("Part 1, naive: ${solveNaive(input.toMutableList())}")
    println("Part 1, unordered: ${solveUnordered(input, count=25)}")
    println("Part 2: ${solveUnordered(input)}")
}
