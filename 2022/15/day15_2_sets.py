#!/usr/bin/env python

import dataclasses
import fileinput
import re
import sys


@dataclasses.dataclass(eq=True, frozen=True)
class Point:
    x: int
    y: int

    def distance(self, other: "Point") -> int:
        return abs(self.x - other.x) + abs(self.y - other.y)


DISTRESS_MIN = 0
DISTRESS_MAX = 4000000


def main():
    line_re = re.compile(
        r"Sensor at x=(?P<sensor_x>[-\d]+), y=(?P<sensor_y>[-\d]+): closest beacon is at x=(?P<beacon_x>[-\d]+), y=(?P<beacon_y>[-\d]+)"
    )
    candidates = {}

    for line in fileinput.input():
        match = line_re.match(line)
        sensor_x = int(match.group("sensor_x"))
        sensor_y = int(match.group("sensor_y"))
        beacon_x = int(match.group("beacon_x"))
        beacon_y = int(match.group("beacon_y"))

        dist = abs(sensor_x - beacon_x) + abs(sensor_y - beacon_y)
        for exclude_y in range(sensor_y - dist, sensor_y + dist + 1):
            if DISTRESS_MIN <= exclude_y <= DISTRESS_MAX:
                delta_x = dist - abs(sensor_y - exclude_y)
                if exclude_y not in candidates:
                    candidates[exclude_y] = set(
                        range(DISTRESS_MIN, DISTRESS_MAX)
                    )
                    size = sys.getsizeof(candidates) + sum(
                        sys.getsizeof(c) for c in candidates
                    )
                    print(size / 1024)
                candidates[exclude_y] -= set(
                    range(sensor_x - delta_x, sensor_x + delta_x + 1)
                )
                if beacon_y == exclude_y:
                    candidates[exclude_y] -= {beacon_x}

    for y in range(DISTRESS_MIN, DISTRESS_MAX + 1):
        if candidates[y]:
            x = list(candidates[y])[0]
            print(x * 4000000 + y)
            break
    else:
        print("no beacon found :(")


if __name__ == "__main__":
    main()
