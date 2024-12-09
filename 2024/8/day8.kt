import kotlin.io.path.Path
import kotlin.io.path.readText

const val DEBUG = false


fun debug(msg: String) {
    if (DEBUG) {
        println(msg)
    }
}


data class Point(val x: Int, val y: Int) {
    operator fun plus(other: Point): Point {
        return Point(x + other.x, y + other.y)
    }

    operator fun minus(other: Point): Point {
        return Point(x - other.x, y - other.y)
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is Point) return false

        return x == other.x && y == other.y
    }

    override fun toString(): String {
        return "(${this.x}, ${this.y})"
    }
}


class GridBounds(val xMax: Int, val yMax: Int) {
    val xMin = 0
    val yMin = 0

    fun contains(loc: Point): Boolean {
        return loc.x >= this.xMin && loc.x < this.xMax &&
            loc.y >= this.yMin && loc.y < this.yMax
    }
}


data class Input(val antennas: Map<Char, List<Point>>,
                 val bounds: GridBounds);


fun readInput(): Input {
    val antennas = mutableMapOf<Char, MutableList<Point>>()
    val lines = Path("input.txt").readText().trim().lines()
    lines.forEachIndexed { y, line ->
        line.forEachIndexed { x, freq ->
            if (freq != '.') {
                antennas.getOrPut(freq) {
                    mutableListOf<Point>()
                }.add(Point(x, y))
            }
        }
    }
    return Input(antennas, GridBounds(lines[0].length, lines.size))
}


fun part1(input: Input): Int {
    val antinodes = mutableListOf<Point>()
    input.antennas.forEach { (freq, antennas) ->
        debug("Collecting antinodes for $freq")
        antennas.forEach { antenna1 ->
            antennas.filter { it != antenna1 }.forEach { antenna2 ->
                debug("Collecting antinodes for $antenna1, $antenna2")
                val delta = antenna1 - antenna2
                debug("  delta: $delta")
                val candidates = listOf(antenna1 + delta, antenna2 - delta)
                debug("  candidates: $candidates")
                antinodes.addAll(candidates.filter { candidate ->
                                     input.bounds.contains(candidate) &&
                                         !antinodes.contains(candidate)
                                 })
                debug("All antinodes: $antinodes")
            }
        }
    }
    return antinodes.size
}


fun part2(input: Input): Int {
    val antinodes = mutableListOf<Point>()
    input.antennas.forEach { (freq, antennas) ->
        debug("Collecting antinodes for $freq")
        antennas.forEach { antenna1 ->
            antennas.filter { it != antenna1 }.forEach { antenna2 ->
                debug("Collecting antinodes for $antenna1, $antenna2")
                val delta = antenna1 - antenna2
                debug("  delta: $delta")
                var cur = antenna1
                while (input.bounds.contains(cur)) {
                    if (!antinodes.contains(cur)) {
                        debug("Found antinode at $cur")
                        antinodes.add(cur)
                    }
                    cur += delta
                }

                cur = antenna2
                while (input.bounds.contains(cur)) {
                    if (!antinodes.contains(cur)) {
                        debug("Found antinode at $cur")
                        antinodes.add(cur)
                    }
                    cur -= delta
                }

                debug("All antinodes: $antinodes")
            }
        }
    }
    return antinodes.size
}


fun main(args: Array<String>) {
    val input = readInput()
    println("Part 1: ${part1(input)}")
    println("Part 2: ${part2(input)}")
}
