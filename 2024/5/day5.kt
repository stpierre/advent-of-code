import kotlin.io.path.Path
import kotlin.io.path.readText

fun debug(msg: String) {
    if (false) {
        println(msg)
    }
}

data class Input(val ordering: Map<String, List<String>>,
                 val updates: List<Update>)

fun readInput(): Input {
    val lines = Path("input.txt").readText().trim().lines()

    val ordering = mutableMapOf<String, MutableList<String>>()
    lines.filter { it.contains("|") }.forEach {
        val (first, second) = it.split("|")
        ordering.getOrPut(first) { mutableListOf<String>() }.add(second)

    }

    val updates = lines.filter { it.contains(",") }.map { Update(it.split(",")) }
    return Input(ordering, updates)
}

class Update(val pages: List<String>) {
    fun isValid(ordering: Map<String, List<String>>): Boolean {
        return this.pages.indices.all { pageIdx ->
            this.pages.slice(pageIdx..this.pages.size - 1).all {
                ordering.get(it)?.contains(this.pages[pageIdx]) != true
            }
        }
    }

    val middlePage = this.pages[(this.pages.size - 1) / 2].toInt()

    override fun toString(): String {
        return this.pages.toString()
    }

    fun reorder(ordering: Map<String, List<String>>): Update {
        val comparator = Comparator<String> { page1, page2 ->
            when {
                ordering.get(page1)?.contains(page2) == true -> -1
                ordering.get(page2)?.contains(page1) == true -> 1
                else -> 0
            }
        }

        return Update(this.pages.sortedWith(comparator))
    }
}

fun part1(ordering: Map<String, List<String>>,
          updates: List<Update>): Int {
    return updates.sumOf { update ->
        if (update.isValid(ordering)) {
            // update is valid, find middle number
            debug("$update is valid, adding ${update.middlePage}")
            update.middlePage
        } else {
            debug("$update is not valid")
            0
        }
    }
}

fun part2(ordering: Map<String, List<String>>,
          updates: List<Update>): Int {
    return updates.sumOf { update ->
        if (!update.isValid(ordering)) {
            debug("$update is not valid")
            val fixedUpdate = update.reorder(ordering)
            debug("Reordered $update -> $fixedUpdate")
            fixedUpdate.middlePage
        } else {
            0
        }
    }
}

fun main(args: Array<String>) {
    val (ordering, updates) = readInput()
    println("Part 1: ${part1(ordering, updates)}")
    println("Part 2: ${part2(ordering, updates)}")
}
