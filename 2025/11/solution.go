package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"maps"
	"os"
	"strings"
	"time"
)

var logger = slog.New(slog.NewTextHandler(
	os.Stderr,
	&slog.HandlerOptions{
		Level: slog.LevelDebug,
		// Level: slog.LevelInfo,
	},
))

type Device struct {
	Label string

	Outputs []*Device
}

func (d Device) String() string {
	return fmt.Sprintf("Device(%s, %d outputs)", d.Label, len(d.Outputs))
}

func (d *Device) dump(dumped map[*Device]struct{}) {
	if _, ok := dumped[d]; ok {
		return
	}
	dumped[d] = struct{}{}

	outputs := make([]string, len((*d).Outputs))
	for i, o := range (*d).Outputs {
		outputs[i] = (*o).Label
		o.dump(dumped)
	}
	fmt.Printf("%s: %s\n", (*d).Label, strings.Join(outputs, " "))
}

func (d Device) Dump() {
	d.dump(make(map[*Device]struct{}))
}

func (d *Device) walk(
	end *Device,
	seen_ map[*Device]struct{},
	cache map[string]int64,
) int64 {
	log := logger.With("from", d.Label, "to", end.Label)
	if d.Label == (*end).Label {
		log.Debug("goal")
		return 1
	}

	if _, ok := seen_[d]; ok {
		log.Debug("cycle detected")
		return 0
	}

	cacheKey := fmt.Sprintf("%s-%s", d.Label, end.Label)
	if result, ok := cache[cacheKey]; ok {
		log.Debug("found cached result", "key", cacheKey, "paths", result)
		return result
	}

	seen := maps.Clone(seen_)
	seen[d] = struct{}{}

	var numPaths int64 = 0
	for _, o := range d.Outputs {
		log.Debug("recursing", "target", o.Label)
		paths := o.walk(end, seen, cache)
		if paths != 0 {
			numPaths += paths
		}
		log.Debug(fmt.Sprintf("%d paths via %s, %d total",
			paths, o.Label, numPaths))
	}

	cache[cacheKey] = numPaths
	return numPaths
}

func (d *Device) Walk(end *Device, c chan int64) {
	retval := d.walk(end, make(map[*Device]struct{}), make(map[string]int64))
	c <- retval
	logger.Info(fmt.Sprintf("%d paths", retval),
		"start", d.Label, "end", end.Label)
}

func (d *Device) AddOutput(target *Device) {
	(*d).Outputs = append((*d).Outputs, target)
}

func part1(start *Device, end *Device) int64 {
	c := make(chan int64)
	go start.Walk(end, c)
	return <-c
}

func part2(svr *Device, dac *Device, fft *Device, out *Device) int64 {
	c := make(chan int64)

	go svr.Walk(dac, c)
	go dac.Walk(fft, c)
	go fft.Walk(out, c)
	var opt1 int64 = 1
	for range 3 {
		opt1 *= <-c
	}
	logger.Info(fmt.Sprintf("Found %d paths svr -> dac -> fft -> out", opt1))

	go svr.Walk(fft, c)
	go fft.Walk(dac, c)
	go dac.Walk(out, c)
	var opt2 int64 = 1
	for range 3 {
		opt2 *= <-c
	}
	logger.Info(fmt.Sprintf("Found %d paths svr -> fft -> dac -> out", opt2))

	return opt1 + opt2
}

func parseLine(line string) (Device, []string) {
	fields := strings.Fields(line)
	device := Device{
		Label:   strings.Trim(fields[0], ":"),
		Outputs: []*Device{},
	}
	return device, fields[1:]
}

func parse() map[string]*Device {
	devices := make(map[string]*Device)
	connections := make(map[string][]string)

	end := &Device{
		Label:   "out",
		Outputs: []*Device{},
	}
	devices[end.Label] = end

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		device, conns := parseLine(strings.TrimSpace(scanner.Text()))
		connections[device.Label] = conns
		devices[device.Label] = &device
	}

	for label, device := range devices {
		for _, conn := range connections[label] {
			target, ok := devices[conn]
			if !ok {
				log.Fatalf("Bad connection: %v -> %s", device, conn)
			}
			device.AddOutput(target)
		}
	}

	return devices
}

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
	devices := parse()
	logger.Info("Parsed input",
		"duration", fmtDuration(time.Now().UnixMicro()-startParsing))

	out := devices["out"]

	var part1Solution int64 = -1
	var duration1 int64
	if you, ok := devices["you"]; ok {
		logger.Info("===== PART 1 =====")
		start1 := time.Now().UnixMicro()
		part1Solution = part1(you, out)
		duration1 = time.Now().UnixMicro() - start1
	} else {
		logger.Error("No 'you' node, skipping part 1")
	}

	var part2Solution int64 = -1
	var duration2 int64
	if svr, ok := devices["svr"]; ok {
		logger.Info("===== PART 2 =====")
		start2 := time.Now().UnixMicro()
		part2Solution = part2(svr, devices["dac"], devices["fft"], out)
		duration2 = time.Now().UnixMicro() - start2
	} else {
		logger.Error("No 'svr' node, skipping part 2")
	}

	if part1Solution != -1 {
		fmt.Printf("Part 1: %d\n", part1Solution)
		fmt.Printf("  in %dμs\n", duration1)
	}
	if part2Solution != -1 {
		fmt.Printf("Part 2: %d\n", part2Solution)
		fmt.Printf("  in %dμs\n", duration2)
	}
}
