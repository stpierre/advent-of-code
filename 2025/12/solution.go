package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

var DEBUG = true

func getLogger() *slog.Logger {
	var level slog.Level
	if DEBUG {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: level},
	))
}

var logger = getLogger()

const (
	PresentCount  = 6
	PresentWidth  = 3
	PresentHeight = 3
)

// ========== PRESENT ==========

type Present [PresentWidth][PresentHeight]bool

func (p Present) Dump() {
	for _, row := range p {
		for _, cell := range row {
			if cell {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func (p Present) Area() int {
	var area int
	for _, row := range p {
		for _, cell := range row {
			if cell {
				area++
			}
		}
	}
	return area
}

type Presents [PresentCount]Present

// ========== TREE ==========

type Tree struct {
	Width        int
	Height       int
	Requirements [PresentCount]int
}

func (t Tree) String() string {
	return fmt.Sprintf("%dx%d %v", t.Width, t.Height, t.Requirements)
}

func (t Tree) TotalPresents() (total int) {
	for _, count := range t.Requirements {
		total += count
	}
	return total
}

func (t Tree) Area() int {
	return t.Width * t.Height
}

func (t Tree) MightFit(p Presents) bool {
	log := logger.With(
		"treeSize", fmt.Sprintf("%dx%d", t.Width, t.Height),
		"treeArea", t.Area(),
		"totalPresents", t.TotalPresents(),
	)
	// this is *not* the same as `t.Area() / 9` because of integer rounding
	maxFit := (t.Width / PresentWidth) * (t.Height / PresentHeight)
	if maxFit >= t.TotalPresents() {
		// presents easily fit in on a 3x3 grid
		log.Debug("tree fits presents trivially", "maxFit", maxFit)
		return true
	}

	return false
}

// ========== SOLVING ==========

func part1(presents Presents, trees []Tree) int {
	count := 0
	for _, t := range trees {
		if t.MightFit(presents) {
			count++
		}
	}
	return count
}

func part2(presents Presents, trees []Tree) int64 {
	return 0
}

// ========== PARSING ==========

func parseTree(line string) Tree {
	fields := strings.Fields(line)

	dims := strings.Trim(fields[0], ":")
	vals := strings.Split(dims, "x")
	width, err := strconv.Atoi(vals[0])
	if err != nil {
		log.Fatalf("Malformed input %s (%s): %v", line, width, err)
	}
	height, err := strconv.Atoi(vals[1])
	if err != nil {
		log.Fatalf("Malformed input %s (%s): %v", line, height, err)
	}

	var reqs [PresentCount]int
	for i, f := range fields[1:] {
		count, err := strconv.Atoi(f)
		if err != nil {
			log.Fatalf("Malformed input %s (%s): %v", line, count, err)
		}
		reqs[i] = count
	}

	return Tree{
		Width:        width,
		Height:       height,
		Requirements: reqs,
	}
}

func parsePresent(lines []string) (p Present) {
	for row, line := range lines {
		for col, char := range line {
			if char == '#' {
				p[row][col] = true
			}
		}
	}
	return p
}

func parse() (presents Presents, trees []Tree) {
	scanner := bufio.NewScanner(os.Stdin)
	inTrees := false
	var presentIdx int
	var presentBuffer []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !inTrees && strings.Contains(line, "x") {
			inTrees = true
		}

		if inTrees {
			trees = append(trees, parseTree(line))
		} else if line == "" {
			presents[presentIdx] = parsePresent(presentBuffer)
			presentBuffer = []string{}
		} else if strings.HasSuffix(line, ":") {
			var err error
			presentIdx, err = strconv.Atoi(strings.Trim(line, ":"))
			if err != nil {
				log.Fatalf("Malformed input %s: %v", line, err)
			}
		} else {
			presentBuffer = append(presentBuffer, line)
		}
	}

	return presents, trees
}

func dump(presents Presents, trees []Tree) {
	for i, present := range presents {
		fmt.Printf("%d:\n", i)
		present.Dump()
		fmt.Print("\n")
	}
	for _, tree := range trees {
		fmt.Println(tree.String())
	}
}

// ========== BOILERPLATE ==========

func fmtDuration(duration_us int64) string {
	if duration_us < 3000 {
		return fmt.Sprintf("%dμs", duration_us)
	} else {
		duration_ms := float64(duration_us) / 1000.0
		if duration_ms < 3000 {
			return fmt.Sprintf("%.2fms", duration_ms)
		} else {
			duration_s := duration_ms / 1000.0
			if duration_s < 180 {
				return fmt.Sprintf("%.2fs", duration_s)
			} else {
				duration_m := duration_s / 60
				if duration_m < 180 {
					return fmt.Sprintf("%.2fm", duration_m)
				} else {
					duration_h := duration_m / 60
					if duration_h < 72 {
						return fmt.Sprintf("%.2fh", duration_h)
					} else {
						duration_d := duration_h / 24
						return fmt.Sprintf("%.2fd", duration_d)
					}
				}
			}
		}
	}
}

func main() {
	logger.Info("===== PARSING =====")
	startParsing := time.Now().UnixMicro()
	presents, trees := parse()
	if DEBUG {
		dump(presents, trees)
	}
	fmt.Printf("Parsed input in %s\n", fmtDuration(time.Now().UnixMicro()-startParsing))

	logger.Info("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(presents, trees)
	duration1 := time.Now().UnixMicro() - start1

	logger.Info("===== PART 2 =====")
	start2 := time.Now().UnixMicro()
	part2Solution := part2(presents, trees)
	duration2 := time.Now().UnixMicro() - start2

	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("  in %dμs\n", duration1)
	fmt.Printf("Part 2: %d\n", part2Solution)
	fmt.Printf("  in %dμs\n", duration2)
}
