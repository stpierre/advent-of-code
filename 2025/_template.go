package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// ========== SETUP ==========

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

// ========== SOLUTION ==========

type PLACEHOLDER error

func part1(data []PLACEHOLDER) int64 {
	return 0
}

func part2(data []PLACEHOLDER) int64 {
	return 0
}

// ========== PARSING ==========

func parseLine(line string) PLACEHOLDER {
	return PLACEHOLDER{}
}

func parse() (data []PLACEHOLDER) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = append(data, parseLine(strings.TrimSpace(scanner.Text())))
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
	data := parse()
	fmt.Printf("Parsed input in %s\n", fmtDuration(time.Now().UnixMicro()-startParsing))

	logger.Info("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(data)
	duration1 := time.Now().UnixMicro() - start1

	logger.Info("===== PART 2 =====")
	start2 := time.Now().UnixMicro()
	part2Solution := part2(data)
	duration2 := time.Now().UnixMicro() - start2

	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("  in %dμs\n", duration1)
	fmt.Printf("Part 2: %d\n", part2Solution)
	fmt.Printf("  in %dμs\n", duration2)
}
