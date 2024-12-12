import kotlin.io.path.Path
import kotlin.io.path.readText

const val DEBUG = false

// input - part 1 expected value - part2 expected value
val TEST_CASES = listOf<Triple<String, Int?, Int?>>(
    Triple("""
AAAA
BBCD
BBCC
EEEC
           """, 140, 80),
    Triple("""
AAAB
ABBA
ABBA
AAAA
           """, 300, 130),
    Triple("""
BAAA
ACAA
AADA
AAAA
           """, 324, 194),
    Triple("""
AAAA
AADA
AADA
AAAA
           """, 320, 120),
    Triple("""
OOOOO
OXOXO
OOOOO
OXOXO
OOOOO
           """, 772, 436),
    Triple("""
OOOOO
OXXXO
OXOXO
OXXXO
OOOOO
           """, 644, 196),
    Triple("""
EEEEE
EXXXX
EEEEE
EXXXX
EEEEE
           """, null, 236),
    Triple("""
AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA
           """, null, 368),
    Triple("""
RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE
           """, 1930, 1206))

fun debug(msg: String) {
    if (DEBUG) {
        println(msg)
    }
}


enum class Direction(val xDelta: Int, val yDelta: Int) {
    NORTH(0, -1),
    EAST(1, 0),
    SOUTH(0, 1),
    WEST(-1, 0);

    fun cw(): Direction {
        val newIdx = (Direction.entries.indexOf(this)!! + 1) % Direction.entries.size
        return Direction.entries[newIdx]
    }

    fun ccw(): Direction {
        val newIdx = (Direction.entries.indexOf(this)!! - 1 + Direction.entries.size) % Direction.entries.size
        return Direction.entries[newIdx]
    }

    fun opposite(): Direction {
        val newIdx = (Direction.entries.indexOf(this)!! + (Direction.entries.size / 2)) % Direction.entries.size
        return Direction.entries[newIdx]
    }
}


data class Point(val x: Int, val y: Int) {
    operator fun plus(other: Point): Point {
        return Point(x + other.x, y + other.y)
    }

    operator fun plus(other: Direction): Point {
        return Point(x + other.xDelta, y + other.yDelta)
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


class Plot(val plantType: Char, val points: List<Point>, val plants: PlantMap) {
    override fun toString(): String {
        return "Plot ${this.plantType} anchored at ${this.points[0]}"
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is Plot) return false

        return this.plantType == other.plantType && this.points == other.points
    }

    fun getArea() = points.size

    private var perimeter: Int? = null

    fun getPerimeter(): Int {
        if (this.perimeter == null) {
            this.perimeter = this.points.sumOf { loc ->
                Direction.entries.map {
                    loc + it
                }.filter {
                    !this.plants.canStep(it, plantType)
                }.count()
            }
            debug("Found perimeter ${this.perimeter} for ${this}")
        }
        return this.perimeter!!
    }

    fun getPart1Cost(): Int {
        val cost = this.getArea() * this.getPerimeter()
        debug("$this costs area ${this.getArea()} * perimeter ${this.getPerimeter()} = $cost")
        return cost
    }

    private var exteriorSides: Int = 0
    private var interiorSides: Int = 0

    fun countSides(exteriorOnly: Boolean = false): Int {
        if (this.exteriorSides == 0) {
            val start = Point(this.points[0]!!.x, this.points[0]!!.y)
            var cur = start
            var direction = Direction.EAST
            val boundaries = mutableListOf<Pair<Point, Direction>>(Pair(cur, direction.opposite()))
            debug("Walking exterior sides of $this starting at $cur going $direction")
            this.exteriorSides = 1
            while (true) {
                val ccw = direction.ccw()
                val nextCcw = cur + ccw
                if (this.plants.canStep(nextCcw, plantType)) {
                    this.exteriorSides++
                    cur = nextCcw
                    direction = ccw
                    debug("${this.plantType}: Turning CCW to $direction, stepping to $cur")
                } else {
                    boundaries.add(Pair(cur, ccw))
                    val next = cur + direction
                    if (this.plants.canStep(next, plantType)) {
                        cur = next
                        debug("${this.plantType}: Continuing $direction, stepping to $cur")
                    } else {
                        boundaries.add(Pair(cur, direction))
                        val cw = direction.cw()
                        val nextCw = cur + cw
                        if (this.plants.canStep(nextCw, plantType)) {
                            this.exteriorSides++
                            cur = nextCw
                            direction = cw
                            debug("${this.plantType}: Turning CW to $direction, stepping to $cur")
                        } else {
                            boundaries.add(Pair(cur, cw))
                            val opp = direction.opposite()
                            val nextOpp = cur + opp
                            if (this.plants.canStep(nextOpp, plantType)) {
                                this.exteriorSides += 2
                                cur = nextOpp
                                direction = opp
                                debug("${this.plantType}: Vault face to $direction, stepping to $cur")
                            } else {
                                boundaries.add(Pair(cur, opp))
                                this.exteriorSides = 4
                                debug("${this.plantType}: No steps possible")
                            }
                        }
                    }
                }
                debug("${this.plantType}: Exterior sides discovered: ${this.exteriorSides}")
                if (cur == start) {
                    if (direction == Direction.WEST) {
                        // we've returned to the starting point from the
                        // same direction we left it. see if we can
                        // continue south, or if we're done.
                        val nextSouth = cur + Direction.SOUTH
                        if (this.plants.canStep(nextSouth, plantType)) {
                            this.exteriorSides++
                            cur = nextSouth
                            direction = Direction.SOUTH
                            debug("${this.plantType}: Turning to $direction from anchor, stepping to $cur")
                        } else {
                            boundaries.add(Pair(cur, Direction.SOUTH))
                            debug("${this.plantType}: Found 1-block tall anchor point at $cur")
                            this.exteriorSides++
                            break
                        }
                    } else {
                        // we've returned to the starting point from
                        // another direction (south), so we're done
                        debug("${this.plantType}: Returned to anchor at $start")
                        break
                    }
                    debug("${this.plantType}: Exterior sides discovered: ${this.exteriorSides}")
                }
            }
            debug("${this.plantType}: Total exterior sides: ${this.exteriorSides}")

            debug("Finding enclaves and interior sides of $this")
            val foundEnclaves = mutableListOf<Plot>()
            this.interiorSides = this.points.sumOf { loc ->
                debug("${this.plantType}: Checking for enclaves from $loc")
                Direction.entries.filter {
                    !boundaries.contains(Pair(loc, it))
                }.map {
                    debug("${this.plantType}: Looking $it from $loc")
                    loc + it
                }.filter {
                    this.plants.contains(it) &&
                        this.plants.getPlantType(it) != this.plantType
                }.map { enclaveStart ->
                    this.plants.plots.find {
                        it.points.contains(enclaveStart)
                    }!!
                }.filter {
                    !foundEnclaves.contains(it)
                }.sumOf { enclave ->
                    foundEnclaves.add(enclave)
                    debug("${this.plantType}: Found enclave $enclave")
                    enclave.countSides(exteriorOnly = true)
                }
            }
            debug("${this.plantType}: Total interior sides discovered: ${this.interiorSides}")
        }

        return this.exteriorSides!! + (if (exteriorOnly) 0 else this.interiorSides!!)
    }

    fun getPart2Cost(): Int {
        val cost = this.getArea() * this.countSides()
        debug("$this costs area ${this.getArea()} * sides ${this.countSides()} = $cost")
        return cost
    }
}


class PlantMap() {
    private val plants = mutableMapOf<Point, Char>()
    private val visited = mutableMapOf<Point, Boolean>()
    val plots = mutableListOf<Plot>()

    fun reset() {
        this.visited.keys.forEach {
            visited[it] = false
        }
    }

    fun add(loc: Point, plantType: Char) {
        this.plants[loc] = plantType
        this.visited[loc] = false
    }

    fun contains(loc: Point): Boolean {
        return this.plants.contains(loc)
    }

    fun findFirstUnvisited(): Point? {
        return this.visited.keys.find {
            !visited[it]!!
        }
    }

    fun getPlantType(loc: Point) = this.plants[loc]!!

    fun isVisited(loc: Point) = this.visited[loc]!!

    fun visit(loc: Point) {
        this.visited[loc] = true
    }

    fun canStep(loc: Point, targetPlantType: Char): Boolean {
        return this.contains(loc) && this.getPlantType(loc) == targetPlantType
    }

    fun getPlot(start: Point): Plot {
        this.visit(start)
        val plantType = this.getPlantType(start)
        val points = mutableSetOf<Point>(start)
        Direction.entries.map {
            start + it
        }.filter {
            this.canStep(it, plantType) && !this.isVisited(it)
        }.forEach {
            points.addAll(this.getPlot(it).points)
        }
        return Plot(plantType, points.toList(), this)
    }

    fun loadAllPlots() {
        while (true) {
            val firstUnvisited = this.findFirstUnvisited()
            if (firstUnvisited == null) {
                debug("No more unvisited plots")
                break
            }

            val plot = this.getPlot(firstUnvisited)
            debug("Got $plot: ${plot.points}")
            this.plots.add(plot)
        }
        debug("Loaded ${this.plots.size} plots")
    }
}


fun readInput(inputText: String? = null): PlantMap {
    val plants = PlantMap()
    val lines = (inputText ?: Path("input.txt").readText()).trim().lines()
    lines.forEachIndexed { y, line ->
        line.forEachIndexed { x, plant ->
            plants.add(Point(x, y), plant)
        }
    }
    return plants
}


fun part1(plants: PlantMap): Int {
    return plants.plots.sumOf { it.getPart1Cost() }
}

fun part2(plants: PlantMap): Int {
    return plants.plots.sumOf { it.getPart2Cost() }
}


fun main(args: Array<String>) {
    TEST_CASES.forEachIndexed { idx, testCase ->
        val (inputText, expected1, expected2) = testCase
        val plants = readInput(inputText)
        plants.loadAllPlots()
        if (expected1 != null) {
            val result1 = part1(plants)
            if (result1 != expected1) {
                error("Got $result1 (!= $expected1) for test case $idx part 1")
            } else {
                println("Test case $idx part 1 passed")
            }
            plants.reset()
        }
        if (expected2 != null) {
            val result2 = part2(plants)
            if (result2 != expected2) {
                error("Got $result2 (!= $expected2) for test case $idx part 2")
            } else {
                println("Test case $idx part 2 passed")
            }
        }
    }

    val plants = readInput()
    plants.loadAllPlots()
    println("Part 1: ${part1(plants)}")
    plants.reset()
    println("Part 2: ${part2(plants)}")
}
