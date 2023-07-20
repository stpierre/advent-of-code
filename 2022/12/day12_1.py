#!/usr/bin/env python

import dataclasses
import fileinput


@dataclasses.dataclass(eq=True, frozen=True)
class Point:
    x: int
    y: int


class Visitor:
    def __init__(self) -> None:
        self.elevations = {}
        self._distances = {}
        self.start = None
        self.target = None
        self._height = 0
        self._width = 0
        self._unvisited = []

    def parse_input(self) -> None:
        for y, line in enumerate(fileinput.input()):
            self._height += 1
            if self._width == 0:
                self._width = len(line.strip())

            for x, char in enumerate(line.strip()):
                current = Point(x, y)
                if char == "S":
                    elevation = 0
                    self.start = current
                    self._distances[current] = 0
                elif char == "E":
                    elevation = 25
                    self.target = current
                    self._distances[current] = None
                else:
                    elevation = ord(char) - 97
                    self._distances[current] = None
                self.elevations[current] = elevation
                self._unvisited.append(current)

    def find_path(self) -> int:
        next_node = self.start
        while next_node is not None:
            self._find_distances(next_node)

            if next_node == self.target:
                return self._distances[self.target]

            next_node = None
            for node in self._unvisited:
                if self._distances[node] is not None and (
                    next_node is None
                    or self._distances[node] < self._distances[next_node]
                ):
                    next_node = node

        raise Exception("No next node")

    def _find_distances(self, current: Point) -> None:
        neighbors = [
            Point(current.x - 1, current.y),
            Point(current.x + 1, current.y),
            Point(current.x, current.y - 1),
            Point(current.x, current.y + 1),
        ]

        for neighbor in neighbors:
            if (
                self._width > neighbor.x >= 0
                and self._height > neighbor.y >= 0
                and neighbor in self._unvisited
                and self.elevations[neighbor] <= self.elevations[current] + 1
            ):
                if (
                    self._distances[neighbor] is None
                    or self._distances[neighbor] > self._distances[current] + 1
                ):
                    self._distances[neighbor] = self._distances[current] + 1
        self._unvisited.remove(current)


def main():
    visitor = Visitor()
    visitor.parse_input()
    print(visitor.find_path())


if __name__ == "__main__":
    main()
