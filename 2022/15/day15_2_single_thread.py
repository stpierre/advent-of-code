#!/usr/bin/env python

import collections
import dataclasses
import fileinput
import re
import time
from typing import Union


@dataclasses.dataclass(eq=True, frozen=True)
class Point:
    x: int
    y: int

    def distance(self, other: "Point") -> int:
        return abs(self.x - other.x) + abs(self.y - other.y)


@dataclasses.dataclass(eq=True, frozen=True)
class XRange:
    start_x: int
    end_x: int
    y: int

    def __contains__(self, other: Union[Point, int]) -> bool:
        if hasattr(other, "y"):
            return self.y == other.y and self.start_x <= other.x <= self.end_x
        return self.start_x <= other <= self.end_x


DISTRESS_MIN = 0
DISTRESS_MAX = 20 # 4000000




def main():
    line_re = re.compile(
        r"Sensor at x=(?P<sensor_x>[-\d]+), y=(?P<sensor_y>[-\d]+): closest beacon is at x=(?P<beacon_x>[-\d]+), y=(?P<beacon_y>[-\d]+)"
    )
    sensors = {}
    beacons = {}
    exclusions = collections.defaultdict(list)
    min_x = max_x = 0
    for line in fileinput.input():
        match = line_re.match(line)
        sensor = Point(
            int(match.group("sensor_x")), int(match.group("sensor_y"))
        )
        beacon = Point(
            int(match.group("beacon_x")), int(match.group("beacon_y"))
        )
        sensors[sensor] = True
        beacons[beacon] = True

        dist = sensor.distance(beacon)
        for exclude_y in range(sensor.y - dist, sensor.y + dist + 1):
            if DISTRESS_MIN <= exclude_y <= DISTRESS_MAX:
                delta_x = dist - abs(sensor.y - exclude_y)
                exclusions[exclude_y].append(
                    XRange(sensor.x - delta_x, sensor.x + delta_x, exclude_y)
                )

        max_x = max(max_x, beacon.x, sensor.x + dist)
        min_x = min(min_x, beacon.x, sensor.x - dist)

    start = time.time()
    for y in range(DISTRESS_MIN, DISTRESS_MAX + 1):
        if y > 0 and y % 1000 == 0:
            elapsed = time.time() - start
            remaining = elapsed * (DISTRESS_MAX / y) / 60 / 60
            print(
                "computed %s rows (%0.2f%%), %0.2fs elapsed, %0.2f hours remaining"
                % (y, y / DISTRESS_MAX, elapsed, remaining)
            )

        candidates = set(range(DISTRESS_MIN, DISTRESS_MAX + 1))
        for beacon in beacons:
            if beacon.y == y:
                candidates -= {beacon.x}

        for exclusion in exclusions.get(y):
            candidates -= set(range(exclusion.start_x, exclusion.end_x + 1))
            if not candidates:
                break

        if len(candidates) > 1:
            raise Exception("oops: %s" % candidates)

        if candidates:
            found = list(candidates)[0]
            print("found beacon at %s, %s" % (found, y))
            print(found * 4000000 + y)
            break


if __name__ == "__main__":
    main()
