#!/usr/bin/env python

import dataclasses
import fileinput


@dataclasses.dataclass(eq=True)
class Point:
    x: int
    y: int

    def move(self, direction: str) -> None:
        if direction == "R":
            self.x += 1
        elif direction == "L":
            self.x -= 1
        elif direction == "D":
            self.y += 1
        elif direction == "U":
            self.y -= 1

    def chase(self, other: "Point") -> None:
        if abs(other.x - self.x) < 2 and abs(other.y - self.y) < 2:
            return

        if other.x != self.x and other.y != self.y:
            # move diagonally
            threshold = 1
        else:
            threshold = 2

        if other.x - self.x >= threshold:
            self.x += 1
        elif other.x - self.x <= -threshold:
            self.x -= 1

        if other.y - self.y >= threshold:
            self.y += 1
        elif other.y - self.y <= -threshold:
            self.y -= 1

    def __str__(self) -> str:
        return f"({self.x}, {self.y})"


def main():
    knots = [Point(0, 0) for _ in range(10)]
    visited = set()
    for line in fileinput.input():
        direction, count = line.strip().split()
        for _ in range(int(count)):
            knots[0].move(direction)
            for i in range(1, len(knots)):
                knots[i].chase(knots[i - 1])
            visited.add(dataclasses.astuple(knots[-1]))
    print(len(visited))


if __name__ == "__main__":
    main()
