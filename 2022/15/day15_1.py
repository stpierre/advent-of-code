#!/usr/bin/env python

import dataclasses
import fileinput
import re
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


TEST_Y = 2000000


def main():
    line_re = re.compile(
        r"Sensor at x=(?P<sensor_x>[-\d]+), y=(?P<sensor_y>[-\d]+): closest beacon is at x=(?P<beacon_x>[-\d]+), y=(?P<beacon_y>[-\d]+)"
    )
    sensors = {}
    beacons = {}
    exclusions = []
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
            if exclude_y == TEST_Y:
                delta_x = dist - abs(sensor.y - exclude_y)
                exclusions.append(
                    XRange(sensor.x - delta_x, sensor.x + delta_x, exclude_y)
                )

        max_x = max(max_x, beacon.x, sensor.x + dist)
        min_x = min(min_x, beacon.x, sensor.x - dist)

    excluded = 0
    for x in range(min_x, max_x + 1):
        test_point = Point(x, TEST_Y)
        if test_point not in beacons and any(
            test_point in e for e in exclusions
        ):
            excluded += 1
    print(excluded)


if __name__ == "__main__":
    main()
