#!/usr/bin/env python

import dataclasses
import fileinput
from concurrent import futures


@dataclasses.dataclass(eq=True, frozen=True)
class Point:
    x: int
    y: int


class NoPathError(Exception):
    """Raised when there is no path from start to target."""


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
                elif char == "E":
                    elevation = 25
                    self.target = current
                else:
                    elevation = ord(char) - 97
                self.elevations[current] = elevation

    @property
    def node_count(self):
        return self._height * self._width

    def find_path(self, start=None) -> int:
        for node in self.elevations.keys():
            self._unvisited.append(node)
            self._distances[node] = None

        next_node = start or self.start
        self._distances[next_node] = 0

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

        raise NoPathError("No next node")

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

    path = visitor.node_count
    procs = []
    with futures.ProcessPoolExecutor() as executor:
        for node, elevation in visitor.elevations.items():
            if elevation == 0:
                procs.append(executor.submit(visitor.find_path, node))

        for future in procs:
            try:
                path_length = future.result()
            except NoPathError:
                continue
            path = min(path, path_length)
    print(path)


if __name__ == "__main__":
    main()
