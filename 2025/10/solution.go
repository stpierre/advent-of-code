package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var DEBUG = false

func Debugf(msg string, subs ...any) {
	if DEBUG {
		log.Printf(msg, subs...)
	}
}

// ========== JOLTAGE ==========

type Joltage []int

func (j Joltage) Increment(b Button) {
	val := b.Value
	for i := 0; val > 0; i++ {
		if val%2 == 1 {
			j[i]++
		}
		val >>= 1
	}
}

func (j Joltage) Decrement(b Button) {
	val := b.Value
	for i := 0; val > 0; i++ {
		if val%2 == 1 {
			j[i]--
		}
		val >>= 1
	}
}

func (j Joltage) Equals(other Joltage) bool {
	for i, v1 := range j {
		if v1 != other[i] {
			return false
		}
	}
	return true
}

// returns true if *any* joltage indicator in j is greater than its counterpart
// in other
func (j Joltage) GreaterThan(other Joltage) bool {
	for i, v1 := range j {
		if v1 > other[i] {
			return true
		}
	}
	return false
}

func (j Joltage) GetCacheKey() string {
	vals := make([]string, len(j))
	for i, v := range j {
		vals[i] = strconv.Itoa(v)
	}
	return strings.Join(vals, ",")
}

func (j Joltage) IsValid() bool {
	for _, val := range j {
		if val < 0 {
			return false
		}
	}
	return true
}

// ========== BUTTONS AND INDICATORS ==========

type Button struct {
	Value  int
	Length int
}

func (b Button) BinaryString() string {
	fmtStr := "%0" + strconv.Itoa(b.Length) + "b"
	return fmt.Sprintf(fmtStr, b.Value)
}

func (b Button) String() string {
	var s []string
	for i := range b.Length {
		if b.Value&int(math.Pow(2, float64(i))) != 0 {
			s = append(s, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("%v", s)
}

type Buttons []Button

func (b Buttons) String() string {
	s := make([]string, len(b))
	for i, button := range b {
		s[i] = button.String()
	}
	return fmt.Sprintf("%v", s)
}

type Indicators struct {
	Value  int
	Length int
}

func (i Indicators) String() string {
	s := make([]string, i.Length)
	for idx := range i.Length {
		if i.Value&int(math.Pow(2, float64(idx))) == 0 {
			s[idx] = "."
		} else {
			s[idx] = "#"
		}
	}
	return strings.Join(s, "")
}

func (i Indicators) BinaryString() string {
	fmtStr := "%0" + strconv.Itoa(i.Length) + "b"
	return fmt.Sprintf(fmtStr, i.Value)
}

// ========== MACHINE ==========

type Machine struct {
	Indicators Indicators
	Buttons    Buttons
	Joltage    Joltage

	length int
	fmtStr string
}

func NewMachine(
	indicators Indicators,
	buttons Buttons,
	joltage Joltage,
) *Machine {
	length := len(joltage)
	fmtStr := "%0" + strconv.Itoa(length) + "b"
	return &Machine{
		Indicators: indicators,
		Buttons:    buttons,
		Joltage:    joltage,
		length:     length,
		fmtStr:     fmtStr,
	}
}

func (m Machine) String() string {
	return fmt.Sprintf("Machine(%v, buttons=%v, joltage=%v)",
		m.Indicators, m.Buttons, m.Joltage)
}

func (m Machine) getNullIndicators() Indicators {
	return Indicators{Value: 0, Length: m.length}
}

func (m Machine) findAllConfigurePaths(
	startState Indicators,
	buttonIdx int,
	goalState Indicators,
) []Buttons {
	prefix := strings.Repeat("  ", buttonIdx)

	b := m.Buttons[buttonIdx]
	if goalState != m.Indicators {
		Debugf("%sConfiguring %v, button %v, starting state=%v, goal state=%v",
			prefix, m.Indicators, b, startState, goalState)
	} else {
		Debugf("%sConfiguring %v, button %v, starting state=%v",
			prefix, m.Indicators, b, startState)
	}

	newState := Indicators{Value: startState.Value ^ b.Value, Length: m.length}
	Debugf("%sPressing button %v: %v => %v", prefix, b, startState, newState)

	var retval []Buttons
	if newState == goalState {
		Debugf("%s-> goal state", prefix)
		retval = append(retval, Buttons{b})
	}

	nextButton := buttonIdx + 1
	if nextButton < len(m.Buttons) {
		for _, c := range m.findAllConfigurePaths(newState, nextButton, goalState) {
			retval = append(retval, append(c, b))
		}

		Debugf("%sNot pressing button %v: %v", prefix, b, startState)
		retval = append(
			retval,
			m.findAllConfigurePaths(startState, nextButton, goalState)...)
	}

	if len(retval) == 0 {
		Debugf("%s-> no solution", prefix)
	}
	return retval
}

func (m Machine) Configure() int {
	Debugf("Configuring %v", m)
	var shortest Buttons
	for _, candidate := range m.findAllConfigurePaths(m.getNullIndicators(), 0, m.Indicators) {
		if len(shortest) == 0 || len(candidate) < len(shortest) {
			shortest = candidate
		}
	}
	Debugf("%v: %d presses: %v", m, len(shortest), shortest)
	return len(shortest)
}

// https://www.reddit.com/r/adventofcode/comments/1pk87hl/2025_day_10_part_2_bifurcate_your_way_to_victory/
func (m Machine) solveJoltage(goal Joltage) Buttons {
	Debugf("Setting joltage for %v => %v", m.Indicators, goal)

	// phase 1: make everything even
	isEven := true
	for _, j := range goal {
		if j%2 != 0 {
			isEven = false
			break
		}
	}
	var options []Buttons
	if !isEven {
		goalIndicators := Indicators{Value: 0, Length: m.length}
		for i, j := range goal {
			if j%2 == 1 {
				goalIndicators.Value += int(math.Pow(2, float64(i)))
			}
		}
		Debugf("%v: To even out %v, goal is %v", m.Indicators, goal, goalIndicators)
		options = m.findAllConfigurePaths(m.getNullIndicators(), 0, goalIndicators)
		if len(options) > 0 {
			Debugf("%v: To even out %v, %d options", m.Indicators, goal, len(options))
		} else {
			Debugf("Cannot even out %v, no solution", m.Indicators)
			return Buttons{}
		}
	} else {
		Debugf("%v: %v is already even", m.Indicators, goal)
		options = append(options, Buttons{})
	}

	var presses Buttons
	for _, opt := range options {
		newGoal := slices.Clone(goal)
		if !isEven {
			for _, b := range opt {
				newGoal.Decrement(b)
			}

			if !newGoal.IsValid() {
				Debugf("%v: Discarding invalid option %v", m.Indicators, opt)
				continue
			}

			Debugf("%v: To even out %v, %d presses %v, new goal=%v",
				m.Indicators, goal, len(opt), opt, newGoal)
		} else {
			Debugf("%v: %v is already even, new goal=%v", m.Indicators, goal, newGoal)
		}
		// phase 2: divide by two and recurse
		success := true
		for i, j := range newGoal {
			if j != 0 {
				success = false
				if j%2 == 0 {
					newGoal[i] = j / 2
				} else {
					log.Fatalf("%v: Goal joltage %d was unexpectedly odd", m, j)
				}
			}
		}

		if success {
			if len(presses) == 0 || len(opt) < len(presses) {
				Debugf("%v: Reached goal in %d", m.Indicators, len(opt))
				presses = opt
			} else {
				Debugf("%v: Reached goal in %d, but already had %d-press path",
					m.Indicators, len(opt), len(presses))
			}
		} else {
			Debugf("%v: Divided goal joltage, recursing on %v",
				m.Indicators, newGoal)
			newPresses := m.solveJoltage(newGoal)
			if len(newPresses) > 0 {
				totalPresses := append(append(newPresses, newPresses...), opt...)
				if len(presses) == 0 || len(totalPresses) < len(presses) {
					Debugf("%v: Recursed to goal in %d", m.Indicators, len(totalPresses))
					presses = totalPresses
				} else {
					Debugf("%v: Recursed to goal in %d, but already had %d-press path",
						m.Indicators, len(totalPresses), len(presses))
				}
			}
		}
	}

	return presses
}

func (m Machine) SetJoltage() int {
	Debugf("Setting joltage for %v", m)
	presses := m.solveJoltage(m.Joltage)
	Debugf("%v: %d presses: %v", m, len(presses), presses)
	return len(presses)
}

func (m Machine) ReplayJoltage(presses Buttons) bool {
	j := make(Joltage, m.length)
	for _, b := range presses {
		j.Increment(b)
	}
	if j.Equals(m.Joltage) {
		fmt.Printf("%v: Verified %d presses\n", m, len(presses))
		return true
	}

	fmt.Printf("%v: FAILED TO VERIFY %d presses: %v\n", m, len(presses), presses)
	return false
}

// ========== SOLUTIONS ==========

func part1(machines []Machine) (totalPresses int) {
	for _, m := range machines {
		totalPresses += m.Configure()
	}
	return totalPresses
}

func part2(machines []Machine) (totalPresses int) {
	for _, m := range machines {
		totalPresses += m.SetJoltage()
	}
	return totalPresses
}

// ========== PARSING ==========

func parseLine(line string) Machine {
	fields := strings.Fields(line)

	var joltage Joltage
	for joltageStr := range strings.SplitSeq(strings.Trim(fields[len(fields)-1], "{}"), ",") {
		j, err := strconv.Atoi(joltageStr)
		if err != nil {
			log.Fatalf("Malformed input %s (%s): %v", line, joltageStr, err)
		}
		joltage = append(joltage, j)
	}

	length := len(joltage)

	indicators := Indicators{Value: 0, Length: length}
	for i, ind := range strings.Trim(fields[0], "[]") {
		length++
		if ind == '#' {
			indicators.Value |= (1 << i)
		}
	}

	var buttons Buttons
	for _, rawButtonSpec := range fields[1 : len(fields)-1] {
		b := Button{Value: 0, Length: length}
		for lightStr := range strings.SplitSeq(strings.Trim(rawButtonSpec, "()"), ",") {
			i, err := strconv.ParseFloat(lightStr, 32)
			if err != nil {
				log.Fatalf("Malformed input %s (%s): %v", line, lightStr, err)
			}
			b.Value |= int(math.Pow(2, i))
		}
		buttons = append(buttons, b)
	}

	m := *NewMachine(indicators, buttons, joltage)
	Debugf("Parsed line %s => %v", line, m)
	return m
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
	Debugf("===== PARSING =====")
	startParsing := time.Now().UnixMicro()
	scanner := bufio.NewScanner(os.Stdin)
	var data []Machine
	for scanner.Scan() {
		data = append(data, parseLine(strings.TrimSpace(scanner.Text())))
	}
	fmt.Printf("Parsed input in %s\n", fmtDuration(time.Now().UnixMicro()-startParsing))

	Debugf("===== PART 1 =====")
	start1 := time.Now().UnixMicro()
	part1Solution := part1(data)
	duration1 := time.Now().UnixMicro() - start1

	Debugf("===== PART 2 =====")
	start2 := time.Now().UnixMicro()
	part2Solution := part2(data)
	duration2 := time.Now().UnixMicro() - start2

	fmt.Printf("Part 1: %d\n", part1Solution)
	fmt.Printf("  in %s\n", fmtDuration(duration1))
	fmt.Printf("Part 2: %d\n", part2Solution)
	fmt.Printf("  in %s\n", fmtDuration(duration2))
}
