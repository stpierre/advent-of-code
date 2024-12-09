import kotlin.io.path.Path
import kotlin.io.path.readText

const val DEBUG = false

fun debug(msg: String) {
    if (DEBUG) {
        println(msg)
    }
}


data class Alloc(var length: Int, var fileId: Long?) {
    override fun toString(): String {
        if (this.fileId == null) {
            return "free(${this.length})"
        } else {
            return "file#${this.fileId}(${this.length})"
        }
    }
}


fun readInput(): List<Alloc> {
    val input = mutableListOf<Alloc>()
    var fileId = 0L
    val filename = if (DEBUG) "test_input.txt" else "input.txt"
    Path(filename).readText().trim().forEachIndexed { idx, char ->
        val length = char.toString().toInt()
        if (idx % 2 == 0) {
            //debug("Read $length blocks of fileID=$fileId")
            input.add(Alloc(length, fileId))
            fileId++
        } else {
            //debug("Read $length free blocks")
            input.add(Alloc(length, null))
        }
        //debug("Input so far: ${dumpAllocs(input)}")
    }
    return input
}


fun dumpBlocks(blocks: List<Long?>): String {
    if (DEBUG) {
        return blocks.map {
            when {
                it == null -> '.'
                else -> it
            }
        }.joinToString("")
    } else {
        return ""
    }
}


fun dumpAllocs(allocs: List<Alloc>, cursor: Int? = null): String {
    if (DEBUG) {
        return allocs.mapIndexed { idx, alloc ->
            val prefix = if (cursor == idx) "\u001b[1m" else ""
            val suffix = if (cursor == idx) "\u001b[0m" else ""
            if (alloc.fileId == null) {
                prefix + ".".repeat(alloc.length) + suffix
            } else {
                prefix + (alloc.fileId.toString() + "-".repeat(alloc.length)).substring(0, alloc.length) + suffix
            }
        }.joinToString("|")
    } else {
        return ""
    }
}


fun checksum(blocks: List<Long?>): Long {
    return blocks.indices.filter {
        blocks[it] != null
    }.sumOf {
        it * blocks[it]!!
    }
}


fun expandAllocs(allocs: List<Alloc>): List<Long?> {
    val blocks = mutableListOf<Long?>()
    allocs.forEach { alloc ->
        if (alloc.fileId != null) {
            //debug("Unpacking $alloc")
            repeat(alloc.length) { blocks.add(alloc.fileId) }
        } else {
            //debug("Unpacking $alloc.length free blocks")
            repeat(alloc.length) { blocks.add(null) }
        }
        //debug("Expanded allocations so far: ${dumpBlocks(blocks)}")
    }
    return blocks
}


fun part1(input: List<Long?>): Long {
    val defragged = input.toMutableList()
    var last = 0
    for (idx in input.size - 1 downTo 0) {
        val fileId = input[idx]
        if (idx <= last) break
        if (fileId == null) continue

        while (last < idx && defragged[last] != null) {
            last++
        }
        if (last < idx) {
            debug("Moving block at $idx containing fileID=$fileId to $last")
            defragged[last] = fileId
            defragged[idx] = null
        } else {
            debug("Finished defragging")
        }
        debug("Blocks: ${dumpBlocks(defragged)}")
    }
    debug("Defragged: ${dumpBlocks(defragged)}")
    return checksum(defragged)
}


fun part2(input: List<Alloc>): Long {
    val defragged = input.toMutableList()
    var cursor = defragged.size
    val moved = mutableListOf<Long>()
    while (true) {
        val maxFree = defragged.slice(0..cursor - 1).filter {
            it.fileId == null
        }.maxByOrNull {
            it.length
        }?.length
        if (maxFree == null) {
            debug("Finished defragging, no more free space")
            break
        }
        val lastMovableIdx = defragged.indices.findLast {
            it <= cursor &&
                defragged[it].fileId != null &&
                !moved.contains(defragged[it].fileId) &&
                defragged[it].length <= maxFree
        }
        if (lastMovableIdx == null) {
            debug("Finished defragging, no allocs to move")
            break
        }
        debug("Cursor at $cursor, maxFree=$maxFree")

        val alloc = defragged[lastMovableIdx]
        val newIdx = defragged.indexOfFirst {
            it.fileId == null &&
                it.length >= alloc.length
        }

        if (newIdx > lastMovableIdx) {
            debug("Cannot move alloc $alloc at $lastMovableIdx to $newIdx")
        } else {
            debug("Moving $alloc at $lastMovableIdx to $newIdx")
            moved.add(alloc.fileId!!)
            val remaining = defragged[newIdx].length - alloc.length
            defragged[newIdx] = alloc

            debug("  Freeing alloc at $lastMovableIdx")
            defragged[lastMovableIdx] = Alloc(alloc.length, null)

            if (remaining > 0) {
                defragged.add(newIdx + 1, Alloc(remaining, null))
                debug("  Inserted alloc for remaining free space at ${newIdx + 1}: ${defragged[newIdx + 1]}")
            }
        }

        cursor = lastMovableIdx
        while (defragged[cursor].fileId == null) {
            cursor--
        }

        debug("Allocations: ${dumpAllocs(defragged, cursor = cursor)}")
    }
    debug("Defragged: ${dumpAllocs(defragged)}")
    return checksum(expandAllocs(defragged))
}


fun main(args: Array<String>) {
    val input = readInput()
    debug("Got input: ${dumpAllocs(input)}")
    println("Part 1: ${part1(expandAllocs(input))}")
    println("Part 2: ${part2(input)}")
}
