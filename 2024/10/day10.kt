import kotlin.io.path.Path
import kotlin.io.path.readText

const val DEBUG = false


fun debug(msg: String) {
    if (DEBUG) {
        println(msg)
    }
}


enum class Direction(val xDelta: Int, val yDelta: Int) {
    NORTH(0, -1),
    SOUTH(0, 1),
    WEST(-1, 0),
    EAST(1, 0),
}


data class Point(val x: Int, val y: Int) {
    operator fun plus(other: Direction): Point {
        return Point(x + other.xDelta, y + other.yDelta)
    }

    override fun toString(): String {
        return "(${this.x}, ${this.y})"
    }
}


class Topography(val heights: List<List<Int>>) {
    val xMin = 0
    val xMax = heights[0]!!.size
    val yMin = 0
    val yMax = heights.size

    fun contains(loc: Point): Boolean {
        return loc.x >= this.xMin && loc.x < this.xMax &&
            loc.y >= this.yMin && loc.y < this.yMax
    }

    fun getHeight(loc: Point): Int {
        if (this.contains(loc)) {
            return this.heights[loc.y][loc.x]
        } else {
            return -1
        }
    }

    fun getTrailheads(): List<Point> {
        val trailheads = mutableListOf<Point>()
        this.heights.indices.forEach { y ->
            this.heights[y].forEachIndexed { x, height ->
                if (height == 0) {
                    trailheads.add(Point(x, y))
                }
            }
        }
        return trailheads
    }

    fun getDestinations(pos: Point, targetHeight: Int = 9): Set<Point> {
        val height = this.getHeight(pos)
        if (height == targetHeight) {
            debug("Successfully ended trail at $pos")
            return setOf(pos)
        }

        val reachable = mutableSetOf<Point>()
        Direction.values().forEach {
            val nextPos = pos + it
            val nextHeight = this.getHeight(nextPos)
            if (nextHeight == height + 1) {
                debug("Walking from $pos ($height) to $nextPos ($nextHeight)")
                reachable.addAll(this.getDestinations(nextPos))
            }
        }
        return reachable
    }

    fun countTrails(pos: Point, targetHeight: Int = 9): Int {
        val height = this.getHeight(pos)
        if (height == targetHeight) {
            debug("Successfully ended trail at $pos")
            return 1
        }

        return Direction.values().sumOf {
            val nextPos = pos + it
            val nextHeight = this.getHeight(nextPos)
            if (nextHeight == height + 1) {
                debug("Walking from $pos ($height) to $nextPos ($nextHeight)")
                this.countTrails(nextPos)
            } else {
                0
            }
        }
    }
}


fun readInput(): Topography {
    val lines = Path("input.txt").readText().trim().lines()
    return Topography(lines.map { line -> line.map { it.toString().toInt() } })
}


fun part1(topo: Topography): Int {
    return topo.getTrailheads().sumOf { trailhead ->
        debug("Trailhead at $trailhead: ${topo.getDestinations(trailhead).size}")
        topo.getDestinations(trailhead).size
    }
}


fun part2(topo: Topography): Int {
    return topo.getTrailheads().sumOf { trailhead ->
        topo.countTrails(trailhead)
    }
}


fun main(args: Array<String>) {
    val input = readInput()
    println("Part 1: ${part1(input)}")
    println("Part 2: ${part2(input)}")
}
