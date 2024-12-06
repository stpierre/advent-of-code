import kotlin.io.path.Path
import kotlin.io.path.readText

const val DEBUG = false

fun readInput(): Input {
    val lines = Path("input.txt").readText().trim().lines()

    val layout = lines.map { line ->
        line.map { it == '#' }
    }

    lines.forEachIndexed { y, line ->
        ICONS.entries.forEach { (direction, icon) ->
            if (line.contains(icon)) {
                return Input(LabMap(layout),
                             Guard(Point(line.indexOf(icon), y), direction))
            }
        }
    }

    error("Guard position not found")
}


fun debug(msg: String) {
    if (DEBUG) {
        println(msg)
    }
}

data class Point(val x: Int, val y: Int) {
    operator fun plus(other: Point): Point {
        return Point(x + other.x, y + other.y)
    }

    operator fun plus(other: Direction): Point {
        return Point(x + other.xDelta, y + other.yDelta)
    }

    override fun toString(): String {
        return "(${this.x}, ${this.y})"
    }
}

enum class Direction(val xDelta: Int, val yDelta: Int) {
    NORTH(0, -1),
    SOUTH(0, 1),
    WEST(-1, 0),
    EAST(1, 0),
}

val ICONS = mapOf(
    Direction.NORTH to '^',
    Direction.SOUTH to 'v',
    Direction.EAST to '>',
    Direction.WEST to '<',
)

fun getDirectionFromIcon(icon: Char): Direction {
    for ((key, value) in ICONS) {
        if (value == icon) {
            return key
        }
    }
    error("No such icon $icon")
}

class LoopException(message: String) : Exception(message)


class LabMap(val layout: List<List<Boolean>>) {
    private val visited = mutableListOf<MutableList<Direction?>>()

    init {
        this.layout.forEach { line ->
            this.visited.add(MutableList(line.size) { null })
        }
    }

    fun isObstructed(loc: Point): Boolean {
        return this.layout[loc.y][loc.x]
    }

    fun visit(guard: Guard) {
        val current = this.visited[guard.location.y][guard.location.x]
        if (current == guard.direction) {
            val msg = "Loop detected: Guard visited ${guard.location} traveling ${guard.direction.name} twice"
            debug(msg)
            throw LoopException(msg)
        }

        debug("Marking ${guard.location} as visited traveling ${guard.direction.name}")
        this.visited[guard.location.y][guard.location.x] = guard.direction
    }

    fun visitedCount(): Int {
        return this.visited.sumOf { row -> row.count { it != null } }
    }

    fun isInBounds(loc: Point): Boolean {
        return (loc.x >= 0 &&
                    loc.x < this.layout[0].size &&
                    loc.y >= 0 &&
                    loc.y < this.layout.size)
    }

    fun dump(guard: Guard): String {
        return this.layout.indices.map { y ->
            this.layout[y].indices.map { x ->
                when {
                    guard.location == Point(x, y) -> guard.icon
                    this.visited[y][x] != null -> ICONS[this.visited[y][x]]
                    this.layout[y][x] -> '#'
                    else -> '.'
                }
            }.joinToString("")
        }.joinToString("\n")
    }

    fun withObstruction(loc: Point): LabMap {
        return LabMap(
            layout.mapIndexed { y, row ->
                row.mapIndexed { x, point ->
                    if (loc == Point(x, y)) {
                        true
                    } else {
                        point
                    }
                }
            }
        )
    }
}

class Guard(var location: Point, var direction: Direction) {
    val icon: Char
        get() = ICONS[this.direction]!!

    fun step(labMap: LabMap): Boolean {
        // returns true if the step remained on the map, false if the
        // guard stepped off the map
        debug("Guard at ${this.location} is taking a step to the ${this.direction.name}")
        var newLoc = this.location + this.direction
        if (!labMap.isInBounds(newLoc)) {
            debug("Guard stepped off map at $newLoc")
            return false
        }
        while (labMap.isObstructed(newLoc)) {
            // if (DEBUG) {
            //     println(labMap.dump(this))
            //     readLine()
            // }
            this.turn()
            debug("New location $newLoc is obstructed, turning to the ${this.direction.name}")
            newLoc = this.location + this.direction
            if (!labMap.isInBounds(newLoc)) {
                debug("Guard stepped off map at $newLoc")
                return false
            }
        }
        this.location = newLoc
        debug("Guard stepped to $newLoc")
        labMap.visit(this)
        return true
    }

    fun turn() {
        if (this.direction == Direction.NORTH) {
            this.direction = Direction.EAST
        } else if (this.direction == Direction.EAST) {
            this.direction = Direction.SOUTH
        } else if (this.direction == Direction.SOUTH) {
            this.direction = Direction.WEST
        } else if (this.direction == Direction.WEST) {
            this.direction = Direction.NORTH
        }
    }

    fun walk(labMap: LabMap) {
        labMap.visit(this)
        while (this.step(labMap)) { }
    }
}

data class Input(val labMap: LabMap, val guard: Guard)

fun part1(labMap: LabMap, guard: Guard): Int {
    guard.walk(labMap)
    if (DEBUG) {
        println(labMap.dump(guard))
    }
    return labMap.visitedCount()
}

fun part2(labMap: LabMap, guard: Guard): Long {
    return labMap.layout.indices.sumOf { y ->
        labMap.layout[y].indices.sumOf { x ->
            val newObstruction = Point(x, y)
            if (!labMap.isObstructed(newObstruction)) {
                debug("Testing map with obstruction at $newObstruction")
                val newMap = labMap.withObstruction(newObstruction)
                val newGuard = Guard(guard.location, guard.direction)
                try {
                    newGuard.walk(newMap)
                    debug("No loop detected with obstruction at $newObstruction")
                    0L
                } catch (e: LoopException) {
                    debug("Loop detected with obstruction at $newObstruction")
                    1L
                }
            }
        }
    }
}

fun main(args: Array<String>) {
    val (labMap, guard) = readInput()
    println("Part 1: ${part1(LabMap(labMap.layout), Guard(guard.location, guard.direction))}")
    println("Part 2: ${part2(LabMap(labMap.layout), Guard(guard.location, guard.direction))}")
}
