import kotlin.io.path.Path
import kotlin.io.path.readText
import kotlin.math.abs

const val DEBUG = false

fun debug(msg: String) {
    if (DEBUG) {
        println(msg)
    }
}

fun Double.isLong(): Boolean {
    return abs(this - this.roundToLong().toDouble()) < 0.001
}

fun Double.roundToLong(): Long {
    return (this + 0.5).toLong()
}

class ClawMachine(val aX: Long, val aY: Long,
                  val bX: Long, val bY: Long,
                  val prizeX: Long, val prizeY: Long) {
    fun solve(): Pair<Long, Long>? {
        val b = ((prizeX - ((aX * prizeY.toDouble()) / aY)) / (bX - ((aX * bY.toDouble()) / aY)))
        debug("$this: B = ${b.toBigDecimal().toPlainString()}")
        if (b.roundToLong() < 0 || !b.isLong()) {
            return null
        }

        val a = ((prizeY - (bY * b)) / aY)
        debug("$this: A = ${a.toBigDecimal().toPlainString()}")
        if (a.roundToLong() < 0 || !a.isLong()) {
            return null
        }
        return Pair(a.roundToLong(), b.roundToLong())
    }

    override fun toString(): String {
        return "Prize: X=${this.prizeX}, Y=${this.prizeY}"
    }
}


fun readInput(): List<ClawMachine> {
    val buttonRE = Regex("Button (.): X\\+(\\d+), Y\\+(\\d+)")
    val prizeRE = Regex("Prize: X=(\\d+), Y=(\\d+)")

    val lines = Path("input.txt").readText().trim().lines()
    return (0..<lines.size step 4).map {
        val buttonA = buttonRE.find(lines[it])!!
        val buttonB = buttonRE.find(lines[it + 1])!!
        val prize = prizeRE.find(lines[it + 2])!!
        ClawMachine(buttonA.groups.get(2)!!.value.toLong(),
                    buttonA.groups.get(3)!!.value.toLong(),
                    buttonB.groups.get(2)!!.value.toLong(),
                    buttonB.groups.get(3)!!.value.toLong(),
                    prize.groups.get(1)!!.value.toLong(),
                    prize.groups.get(2)!!.value.toLong())
    }
}


fun part1(input: List<ClawMachine>): Long {
    return input.sumOf {
        val presses = it.solve()
        if (presses == null) {
            debug("No solution for $it")
            0L
        } else {
            val (pressA, pressB) = presses!!
            debug("Solved $it: press A $pressA times, press B $pressB times")
            pressA * 3 + pressB
        }
    }
}


fun part2(input: List<ClawMachine>): Long {
    val adj = 10000000000000L
    return input.map {
        ClawMachine(it.aX, it.aY, it.bX, it.bY,
                    it.prizeX + adj, it.prizeY + adj)
    }.sumOf {
        val presses = it.solve()
        if (presses == null) {
            debug("No solution for $it")
            0L
        } else {
            val (pressA, pressB) = presses!!
            debug("Solved $it: press A $pressA times, press B $pressB times")
            pressA * 3 + pressB
        }
    }
}

fun main(args: Array<String>) {
    val input = readInput()
    println("Part 1: ${part1(input)}")
    println("Part 2: ${part2(input)}")
}
