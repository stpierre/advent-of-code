package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var DEBUG = true

func Debugf(msg string, subs ...any) {
	if DEBUG {
		log.Printf(msg, subs...)
	}
}

type Point struct {
	X int
	Y int
	Z int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Z)
}

func (p Point) Distance(other Point) float64 {
	return math.Sqrt(
		math.Pow(float64(other.X-p.X), 2) +
			math.Pow(float64(other.Y-p.Y), 2) +
			math.Pow(float64(other.Z-p.Z), 2))
}

type Connection struct {
	Points   [2]Point
	Distance float64
}

func (c Connection) String() string {
	return fmt.Sprintf("%s:%s", c.Points[0].String(), c.Points[1].String())
}

func NewConnection(p1 Point, p2 Point) *Connection {
	dist := p1.Distance(p2)
	return &Connection{
		Points:   [2]Point{p1, p2},
		Distance: dist,
	}
}

type Circuit map[Point]struct{}

func (c Circuit) ContainsEither(conn Connection) bool {
	_, ok1 := c[conn.Points[0]]
	if ok1 {
		return true
	}
	_, ok2 := c[conn.Points[1]]
	return ok2
}

func (c Circuit) Dump() {
	if DEBUG {
		Debugf("Circuit contains %d points:", len(c))
		for p := range c {
			Debugf("  %s", p)
		}
	}
}

func (c *Circuit) AddPoint(p Point) {
	(*c)[p] = struct{}{}
}

func (c *Circuit) AddConnection(conn Connection) {
	for _, p := range conn.Points {
		c.AddPoint(p)
	}
}

func (c *Circuit) Merge(c2 Circuit) {
	for p := range c2 {
		(*c)[p] = struct{}{}
	}
}

func NewCircuit(conn Connection) *Circuit {
	c := make(Circuit)
	c.AddConnection(conn)
	return &c
}

func getPotentialConnections(data []Point) (result []Connection) {
	for i, p1 := range data {
		for j := i + 1; j < len(data); j++ {
			result = append(result, *NewConnection(p1, data[j]))
		}
	}
	return result
}

func connectN(data []Point, count int) []Circuit {
	conns := getPotentialConnections(data)
	sort.Slice(conns, func(i, j int) bool {
		return conns[i].Distance < conns[j].Distance
	})
	var circuits []Circuit
	for i, conn := range conns {
		if i >= count {
			Debugf("%d connections made, breaking out", count)
			break
		}

		var containedIn []int
		for i, circuit := range circuits {
			// circuit.Dump()
			if circuit.ContainsEither(conn) {
				containedIn = append(containedIn, i)
			}
		}

		if len(containedIn) == 2 {
			Debugf("Connection %s bridges two circuits, merging", conn.String())
			circuits[containedIn[0]].Merge(circuits[containedIn[1]])
			// remove the extra circuit
			circuits = append(circuits[:containedIn[1]], circuits[containedIn[1]+1:]...)
		} else if len(containedIn) == 1 {
			Debugf("Adding connection %s to existing circuit", conn.String())
			circuits[containedIn[0]].AddConnection(conn)
		} else {
			Debugf("Creating new circuit for %s", conn.String())
			c := make(Circuit)
			c.AddConnection(conn)
			circuits = append(circuits, c)
		}
	}
	Debugf("Returning %d circuits", len(circuits))
	return circuits
}

func part1(data []Point) int {
	circuits := connectN(data, 1000)
	sort.Slice(circuits, func(i, j int) bool {
		return len(circuits[i]) > len(circuits[j])
	})
	Debugf("Three largest circuits: %d, %d, %d", len(circuits[0]), len(circuits[1]), len(circuits[2]))
	return len(circuits[0]) * len(circuits[1]) * len(circuits[2])
}

func part2(data []Point) int {
	conns := getPotentialConnections(data)
	sort.Slice(conns, func(i, j int) bool {
		return conns[i].Distance < conns[j].Distance
	})
	var circuits []Circuit
	for _, p := range data {
		c := make(Circuit)
		c.AddPoint(p)
		circuits = append(circuits, c)
	}

	for i, conn := range conns {
		var containedIn []int
		for i, circuit := range circuits {
			// circuit.Dump()
			if circuit.ContainsEither(conn) {
				containedIn = append(containedIn, i)
			}
		}

		if len(containedIn) == 2 {
			Debugf("Connection %s bridges two circuits, merging", conn.String())
			circuits[containedIn[0]].Merge(circuits[containedIn[1]])
			// remove the extra circuit
			circuits = append(circuits[:containedIn[1]], circuits[containedIn[1]+1:]...)
		} else if len(containedIn) == 1 {
			Debugf("Adding connection %s to existing circuit", conn.String())
			circuits[containedIn[0]].AddConnection(conn)
		} else {
			Debugf("Creating new circuit for %s", conn.String())
			c := make(Circuit)
			c.AddConnection(conn)
			circuits = append(circuits, c)
		}
		if len(circuits) == 1 {
			Debugf("%d connections made, final: %s", i, conn.String())
			return conn.Points[0].X * conn.Points[1].X
		}
	}

	log.Fatalf("%d circuits remaining after all connections made", len(circuits))
	return 0
}

func parseLine(line string) Point {
	vals := strings.Split(line, ",")
	x, err := strconv.Atoi(vals[0])
	if err != nil {
		log.Fatalf("Malformed input %s (%s): %v", line, x, err)
	}

	y, err := strconv.Atoi(vals[1])
	if err != nil {
		log.Fatalf("Malformed input %s (%s): %v", line, y, err)
	}

	z, err := strconv.Atoi(vals[2])
	if err != nil {
		log.Fatalf("Malformed input %s (%s): %v", line, z, err)
	}
	return Point{X: x, Y: y, Z: z}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var data []Point
	for scanner.Scan() {
		data = append(data, parseLine(strings.TrimSpace(scanner.Text())))
	}

	Debugf("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(data)
	duration1 := time.Now().UnixMicro() - start1

	Debugf("===== PART 2 =====")
	start2 := time.Now().UnixMicro()
	part2Solution := part2(data)
	duration2 := time.Now().UnixMicro() - start2

	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("  in %dμs\n", duration1)
	fmt.Printf("Part 2: %d\n", part2Solution)
	fmt.Printf("  in %dμs\n", duration2)
}
