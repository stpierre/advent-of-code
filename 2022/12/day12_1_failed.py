#!/usr/bin/env python

import copy
import dataclasses
import fileinput
import random
import sys


@dataclasses.dataclass(eq=True, frozen=True)
class Point:
    x: int
    y: int


class NavigationError(Exception):
    pass


DEAD_ENDS = []
EXPLORED = set()
BEST_PATH = None


def traverse(elevations, current, target, traversed=None):
    global BEST_PATH
    traversed = traversed or []
    if BEST_PATH and len(traversed) + 1 > len(BEST_PATH):
        raise NavigationError("too long")

    if current in traversed:
        raise NavigationError("loop")

    if current == target:
        return traversed

    traversed.append(current)

    EXPLORED.add(current)
    visited = len(EXPLORED) / (len(elevations) * len(elevations[0]))
    if random.random() < 0.01:
        if not BEST_PATH:
            print(
                "%s%% visited, found %s dead ends, no path found yet"
                % (visited * 100, len(DEAD_ENDS))
            )
        else:
            print(
                "%s%% visited, found %s dead ends, best path = %s"
                % (visited * 100, len(DEAD_ENDS), len(BEST_PATH))
            )

    current_elevation = elevations[current.y][current.x]

    options = [
        Point(current.x - 1, current.y),
        Point(current.x + 1, current.y),
        Point(current.x, current.y - 1),
        Point(current.x, current.y + 1),
    ]

    paths = []
    for option in options:
        if (
            option.x >= 0
            and option.x < len(elevations[0])
            and option.y >= 0
            and option.y < len(elevations)
            and option not in DEAD_ENDS
            and option not in traversed
            and elevations[option.y][option.x] <= current_elevation + 1
        ):
            try:
                paths.append(
                    traverse(
                        elevations,
                        option,
                        target,
                        traversed=copy.copy(traversed),
                    )
                )
            except NavigationError:
                pass
    if not paths:
        DEAD_ENDS.append(current)
        raise NavigationError("inescapable gorge")
    best = min(paths, key=len)
    if not BEST_PATH or len(best) < len(BEST_PATH):
        BEST_PATH = best
    return best


def main():
    elevations = []
    start = None
    target = None

    for y, line in enumerate(fileinput.input()):
        row = []
        for x, char in enumerate(line.strip()):
            if char == "S":
                row.append(0)
                start = Point(x, y)
            elif char == "E":
                row.append(25)
                target = Point(x, y)
            else:
                row.append(ord(char) - 97)
        elevations.append(row)

    sys.setrecursionlimit(len(elevations) * len(elevations[0]))

    steps = traverse(elevations, start, target)
    print(len(steps))


if __name__ == "__main__":
    main()
