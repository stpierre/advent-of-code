package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
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

// ========== DIRECTIONS ==========

type DirectionFamily string

const (
	Vertical   DirectionFamily = "vertical"
	Horizontal DirectionFamily = "horizontal"
)

type Direction struct {
	Point
	Family DirectionFamily
	Name   string
}

func (d Direction) String() string {
	return d.Name
}

var (
	Up    = Direction{Name: "up", Point: Point{X: 0, Y: -1}, Family: Vertical}
	Right = Direction{Name: "right", Point: Point{X: 1, Y: 0}, Family: Horizontal}
	Down  = Direction{Name: "down", Point: Point{X: 0, Y: 1}, Family: Vertical}
	Left  = Direction{Name: "left", Point: Point{X: -1, Y: 0}, Family: Horizontal}
)

// ========== POINT ==========

type Point struct {
	X int
	Y int
}

func (p Point) isOn(e Edge) bool {
	return e.Contains(p)
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

// ========== EDGE ==========

type Edge struct {
	Start Point
	End   Point

	Direction Direction
}

func NewEdge(start Point, end Point) *Edge {
	var dir Direction
	if start.X == end.X {
		if start.Y > end.Y {
			dir = Down
		} else {
			dir = Up
		}
	} else {
		if start.X > end.X {
			dir = Left
		} else {
			dir = Right
		}
	}

	return &Edge{
		Start:     start,
		End:       end,
		Direction: dir,
	}
}

func (e Edge) Range() (Point, Point) {
	if e.Start.X < e.End.X || e.Start.Y < e.End.Y {
		return e.Start, e.End
	}
	return e.End, e.Start
}

func (e Edge) Min() Point {
	retval, _ := e.Range()
	return retval
}

func (e Edge) Max() Point {
	_, retval := e.Range()
	return retval
}

func (e Edge) String() string {
	return fmt.Sprintf("%v(%v-%v)", e.Direction, e.Start, e.End)
}

func (e Edge) IterPoints(yield func(Point) bool) {
	start, end := e.Range()
	for x := start.X; x <= end.X; x++ {
		for y := end.Y; y <= end.Y; y++ {
			if !yield(Point{X: x, Y: y}) {
				return
			}
		}
	}
}

func (e Edge) Contains(p Point) bool {
	start, end := e.Range()
	if e.Direction.Family == Vertical {
		return start.X == p.X && p.Y >= start.Y && p.Y <= end.Y
	}
	// else, horizontal edge
	return start.Y == p.Y && p.X >= start.X && p.X <= end.X
}

func (e Edge) Intersects(other Edge) bool {
	if e.Direction.Family == other.Direction.Family {
		return other.Contains(e.Start) || other.Contains(e.End) ||
			e.Contains(other.Start) || e.Contains(other.End)
	}

	var horiz Edge
	var vert Edge
	if e.Direction.Family == Horizontal {
		horiz, vert = e, other
	} else {
		horiz, vert = other, e
	}
	hMin, hMax := horiz.Range()
	vMin, vMax := vert.Range()
	return vMin.X > hMin.X && vMax.X < hMax.X &&
		hMin.Y > vMin.Y && hMax.Y < vMax.Y
}

// ========== PERIMETER ==========

type Perimeter struct {
	Edges []Edge
	Min   Point
	Max   Point
}

func NewPerimeter(vertices []Point) *Perimeter {
	xMin := vertices[0].X
	xMax := 0
	yMin := vertices[0].Y
	yMax := 0

	var edges []Edge
	prev := vertices[len(vertices)-1]
	for _, p := range vertices {
		xMin = min(xMin, p.X)
		xMax = max(xMax, p.X)
		yMin = min(yMin, p.Y)
		yMax = max(yMax, p.Y)

		Debugf("Constructing new edge %v - %v", prev, p)
		edges = append(edges, *NewEdge(prev, p))
		prev = p
	}

	return &Perimeter{
		Edges: edges,
		Min:   Point{X: xMin, Y: yMin},
		Max:   Point{X: xMax, Y: yMax},
	}
}

func (p *Perimeter) GetVertices() []Point {
	var retval []Point
	for _, e := range p.Edges {
		retval = append(retval, e.Start)
	}
	Debugf("Got %d vertices", len(retval))
	return retval
}

func (p *Perimeter) Contains(a Area) bool {
	// not a general-purpose solution
	bounds := make(map[string]bool)
	for _, edge := range p.Edges {
		eMin, eMax := edge.Range()
		if edge.Direction.Family == Horizontal {
			if eMin.X <= a.Left && eMax.X >= a.Left {
				if eMin.Y >= a.Top {
					bounds["above top left"] = true
				} else if eMin.Y <= a.Bottom {
					bounds["below bottom left"] = true
				} else if eMax.X > a.Left {
					Debugf("%v is intersected by %v", a, edge)
					return false
				}
			}
			if eMin.X <= a.Right && eMax.X >= a.Right {
				if eMin.Y >= a.Top {
					bounds["above top right"] = true
				} else if eMin.Y <= a.Bottom {
					bounds["below bottom right"] = true
				} else if eMin.X < a.Right {
					Debugf("%v is intersected by %v", a, edge)
					return false
				}
			}
		} else { // Vertical edge
			if eMin.Y <= a.Top && eMax.Y >= a.Top {
				if eMin.X >= a.Right {
					bounds["right of top right"] = true
				} else if eMin.X <= a.Left {
					bounds["left of top left"] = true
				} else if eMin.Y < a.Top {
					Debugf("%v is intersected by %v", a, edge)
					return false
				}
			}
			if eMin.Y <= a.Bottom && eMax.Y >= a.Bottom {
				if eMin.X >= a.Right {
					bounds["right of bottom right"] = true
				} else if eMin.X <= a.Left {
					bounds["left of bottom left"] = true
				} else if eMax.Y > a.Bottom {
					Debugf("%v is intersected by %v", a, edge)
					return false
				}
			}
		}

		if len(bounds) == 8 {
			Debugf("%v is bounded on all sides", a)
			return true
		}
	}
	Debugf("%v is bounded on %d sides: %v", a, len(bounds), maps.Keys(bounds))
	return false
}

func (p Perimeter) AssertNoEdgeLoops() {
	Debugf("Checking for overlapping edges")
	var horizontal []Edge
	var vertical []Edge
	for _, e := range p.Edges {
		if e.Direction.Family == Horizontal {
			horizontal = append(horizontal, e)
		} else {
			vertical = append(vertical, e)
		}
	}

	start := time.Now().UnixMicro()
	for checked, horiz := range horizontal {
		if checked > 0 && checked%100 == 0 {
			elapsed_us := float64(time.Now().UnixMicro() - start)
			complete := float64(checked) / float64(len(p.Edges))
			expected_us := int(((1 / complete) * elapsed_us) - elapsed_us)
			Debugf("Checked %d/%d edges in %dμs, %dμs remaining",
				checked, len(p.Edges), int(elapsed_us), expected_us)
		}

		hMin, hMax := horiz.Range()
		for _, vert := range vertical {
			vMin, vMax := vert.Range()
			if vMin.X > hMin.X && vMax.X < hMax.X &&
				hMin.Y > vMin.Y && hMax.Y < vMax.Y {
				log.Fatalf("Found overlapping edges: %v, %v", horiz, vert)
			}
		}
	}
}

// ========== AREAS ==========

type Area struct {
	Top    int
	Left   int
	Bottom int
	Right  int
}

func NewArea(p1 Point, p2 Point) *Area {
	return &Area{
		Top:    max(p1.Y, p2.Y),
		Left:   min(p1.X, p2.X),
		Bottom: min(p1.Y, p2.Y),
		Right:  max(p1.X, p2.X),
	}
}

func (a Area) TopLeft() Point {
	return Point{X: a.Left, Y: a.Top}
}

func (a Area) TopRight() Point {
	return Point{X: a.Right, Y: a.Top}
}

func (a Area) BottomRight() Point {
	return Point{X: a.Right, Y: a.Bottom}
}

func (a Area) BottomLeft() Point {
	return Point{X: a.Left, Y: a.Bottom}
}

func (a Area) Vertices() [4]Point {
	return [4]Point{a.TopLeft(), a.TopRight(), a.BottomRight(), a.BottomLeft()}
}

func (a Area) String() string {
	return fmt.Sprintf("Area(%v - %v)", a.TopLeft(), a.BottomRight())
}

func (a Area) Area() int {
	h := float64(a.Right-a.Left) + 1
	w := float64(a.Top-a.Bottom) + 1
	return int(h * w)
}

// ...and the rest

func part1(p Perimeter) (maxArea int) {
	v := p.GetVertices()
	for i, p1 := range v {
		for j := i; j < len(v); j++ {
			a := NewArea(p1, v[j])
			maxArea = max(maxArea, a.Area())
		}
	}
	return maxArea
}

func part2(p Perimeter) (maxArea int) {
	v := p.GetVertices()
	Debugf("Vertices: %v", v)
	for i, p1 := range v {
		Debugf("p1 = %v", p1)
		for j := i; j < len(v); j++ {
			Debugf("p1=%v, i=%d, j=%d, p2=%v", p1, i, j, v[j])
			a := *NewArea(p1, v[j])
			Debugf("Checking %v", a)
			if p.Contains(a) {
				maxArea = max(maxArea, a.Area())
			}
		}
	}
	return maxArea
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

	return Point{X: x, Y: y}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var vertices []Point
	for scanner.Scan() {
		vertices = append(vertices, parseLine(strings.TrimSpace(scanner.Text())))
	}
	p := *NewPerimeter(vertices)

	// there are no loops in the edges
	p.AssertNoEdgeLoops()

	Debugf("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(p)
	duration1 := time.Now().UnixMicro() - start1

	Debugf("===== PART 2 =====")
	start2 := time.Now().UnixMicro()
	part2Solution := part2(p)
	duration2 := time.Now().UnixMicro() - start2

	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("  in %dμs\n", duration1)
	fmt.Printf("Part 2: %d\n", part2Solution)
	fmt.Printf("  in %dμs\n", duration2)
}
