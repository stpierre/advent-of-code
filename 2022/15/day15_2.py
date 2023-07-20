#!/usr/bin/env python

import collections
import dataclasses
import fileinput
import re
import time
from concurrent import futures
from typing import Union


@dataclasses.dataclass(eq=True, frozen=True)
class Point:
    x: int
    y: int

    def distance(self, other: "Point") -> int:
        return abs(self.x - other.x) + abs(self.y - other.y)


class UnmergeableError(Exception):
    """Raised when trying to merge two XRanges that cannot be merged."""


@dataclasses.dataclass(eq=True, frozen=True)
class XRange:
    start_x: int
    end_x: int
    y: int

    def __contains__(self, other: Union[Point, int]) -> bool:
        if hasattr(other, "y"):
            return self.y == other.y and self.start_x <= other.x <= self.end_x
        return self.start_x <= other <= self.end_x

    def merge(self, other: "XRange") -> "XRange":
        if self.y != other.y:
            raise ValueError("y coordinates do not match")

        if (
            self.start_x == other.end_x + 1
            or other.start_x == self.end_x + 1
            or other.start_x <= self.end_x <= other.end_x
            or self.start_x <= other.end_x <= self.end_x
        ):
            return XRange(
                min(self.start_x, other.start_x),
                max(self.end_x, other.end_x),
                self.y,
            )

        raise UnmergeableError(f"Cannot merge {self} and {other}")

    def __str__(self):
        return f"({self.start_x}, {self.y}) to ({self.end_x}, {self.y})"


class Locator:
    line_re = re.compile(
        r"Sensor at x=(?P<sensor_x>[-\d]+), y=(?P<sensor_y>[-\d]+): closest beacon is at x=(?P<beacon_x>[-\d]+), y=(?P<beacon_y>[-\d]+)"
    )

    def __init__(self, distress_max):
        self.distress_min = 0
        self.distress_max = distress_max
        self.beacons = {}
        self.sensors = {}
        self.exclusions = collections.defaultdict(list)

    def find_beacon(self, y):
        candidates = set(range(self.distress_min, self.distress_max + 1))
        for beacon in self.beacons:
            if beacon.y == y:
                candidates -= {beacon.x}

        for exclusion in self.exclusions.pop(y):
            candidates -= set(range(exclusion.start_x, exclusion.end_x + 1))
            if not candidates:
                return None

        if len(candidates) > 1:
            raise Exception("oops: %s" % candidates)
        return list(candidates)[0], y

    def parse_input(self):
        for line in fileinput.input():
            match = self.line_re.match(line)
            sensor = Point(
                int(match.group("sensor_x")), int(match.group("sensor_y"))
            )
            beacon = Point(
                int(match.group("beacon_x")), int(match.group("beacon_y"))
            )
            self.beacons[beacon] = True

            dist = sensor.distance(beacon)
            for exclude_y in range(sensor.y - dist, sensor.y + dist + 1):
                if self.distress_min <= exclude_y <= self.distress_max:
                    delta_x = dist - abs(sensor.y - exclude_y)
                    x_range = XRange(
                        sensor.x - delta_x, sensor.x + delta_x, exclude_y
                    )
                    for exclusion in self.exclusions[exclude_y]:
                        try:
                            merged = exclusion.merge(x_range)
                        except UnmergeableError:
                            pass
                        else:
                            self.exclusions[exclude_y].remove(exclusion)
                            self.exclusions[exclude_y].append(merged)
                            break
                    else:
                        self.exclusions[exclude_y].append(x_range)

    def locate(self):
        start = time.time()
        procs = []
        with futures.ProcessPoolExecutor() as executor:
            for y in range(self.distress_min, self.distress_max + 1):
                procs.append(executor.submit(self.find_beacon, y))

            for i, future in enumerate(procs):
                found = future.result()
                if i > 0 and i % 1000 == 0:
                    elapsed = time.time() - start
                    remaining = elapsed * (self.distress_max / y) / 60 / 60
                    print(
                        "computed %s rows (%0.2f%%), %0.2fs elapsed, %0.2f hours remaining"
                        % (i, i / self.distress_max, elapsed, remaining)
                    )

                if found:
                    found_x, found_y = found
                    return found_x * 4000000 + found_y


def main():
    loc = Locator(4000000)
    loc.parse_input()
    import sys

    print(sys.getsizeof(loc.exclusions))
    # print(loc.locate())


if __name__ == "__main__":
    main()
