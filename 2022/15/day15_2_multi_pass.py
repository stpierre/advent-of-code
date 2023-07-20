#!/usr/bin/env python

import fileinput
import operator
import re
import time

DISTRESS_MIN = 0
DISTRESS_MAX = 4000000


def main():
    line_re = re.compile(
        r"Sensor at x=(?P<sensor_x>[-\d]+), y=(?P<sensor_y>[-\d]+): closest beacon is at x=(?P<beacon_x>[-\d]+), y=(?P<beacon_y>[-\d]+)"
    )
    sensors = []

    for line in fileinput.input():
        match = line_re.match(line)
        sensor_x = int(match.group("sensor_x"))
        sensor_y = int(match.group("sensor_y"))
        beacon_x = int(match.group("beacon_x"))
        beacon_y = int(match.group("beacon_y"))
        dist = abs(sensor_x - beacon_x) + abs(sensor_y - beacon_y)
        sensors.append((sensor_x, sensor_y, beacon_x, beacon_y, dist))

    done = False
    start = time.time()
    for y in range(DISTRESS_MIN, DISTRESS_MAX + 1):
        ranges = []
        for sensor_x, sensor_y, beacon_x, beacon_y, dist in sensors:
            if sensor_y - dist <= y <= sensor_y + dist:
                delta_x = dist - abs(sensor_y - y)
                ranges.append((sensor_x - delta_x, sensor_x + delta_x))
                if beacon_y == y:
                    ranges.append((beacon_x, beacon_x))
        ranges.sort(key=operator.itemgetter(0))

        if ranges[0][0] > DISTRESS_MIN:
            print(y)
            break
        else:
            exclude_end = 0
            for x_range in ranges:
                if x_range[0] > exclude_end + 1:
                    x = exclude_end + 1
                    print(x * 4000000 + y)
                    done = True
                    break
                else:
                    exclude_end = max(exclude_end, x_range[1])
            else:
                if exclude_end < DISTRESS_MAX:
                    print(DISTRESS_MAX * 4000000 + y)
                    break

        if done:
            break


if __name__ == "__main__":
    main()
