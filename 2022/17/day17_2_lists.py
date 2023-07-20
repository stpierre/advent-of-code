#!/usr/bin/env python

import dataclasses
import fileinput
import itertools
import time
from typing import List, Tuple

WIDTH = 7
ROCKS = 2022
# ROCKS = 1000000000000


class Shape:
    height = 0
    width = 0

    def __init__(self, bottom: int, left: int = 2):
        self.bottom = bottom
        self.left = left

    @property
    def right(self):
        return self.left + self.width - 1

    @property
    def top(self):
        return self.bottom + self.height - 1

    def can_blow(
        self, direction: str, occupied: List[Tuple[int, int]]
    ) -> bool:
        right = self.right
        if direction == ">":
            return right < WIDTH - 1 and not any(
                [right + 1, y] in occupied
                for y in range(self.bottom, self.top + 1)
            )
        if direction == "<":
            return self.left > 0 and not any(
                [self.left - 1, y] in occupied
                for y in range(self.bottom, self.top + 1)
            )

        raise Exception(f"Unknown direction {direction}")

    def blow(self, direction: str, occupied: List[Tuple[int, int]]) -> None:
        if not self.can_blow(direction, occupied):
            return
        if direction == ">":
            self.left += 1
        elif direction == "<":
            self.left -= 1
        else:
            raise Exception(f"Unknown direction {direction}")

    def can_drop(self, occupied: List[Tuple[int, int]]) -> bool:
        return self.bottom - 1 > 0 and not any(
            [x, self.bottom - 1] in occupied
            for x in range(self.left, self.right + 1)
        )

    def drop(self, heights: List[int]) -> Tuple[bool, List[int]]:
        if not self.can_drop(heights):
            return True, self.freeze()
        self.bottom -= 1
        return False, heights

    def freeze(self) -> List[Tuple[int, int]]:
        retval = []
        for x in range(self.left, self.right + 1):
            for y in range(self.bottom, self.top + 1):
                retval.append([x, y])
        return retval


class Flat(Shape):
    height = 1
    width = 4


class Plus(Shape):
    height = 3
    width = 3

    def can_blow(
        self, direction: str, occupied: List[Tuple[int, int]]
    ) -> bool:
        if direction == ">":
            right = self.right
            return (
                right < WIDTH - 1
                and [right + 1, self.bottom + 1] not in occupied
                and [right, self.bottom] not in occupied
                and [right, self.top] not in occupied
            )

        if direction == "<":
            return (
                self.left > 0
                and [self.left - 1, self.bottom + 1] not in occupied
                and [self.left, self.bottom] not in occupied
                and [self.left, self.top] not in occupied
            )
        return super().can_blow(direction, occupied)

    def can_drop(self, occupied: List[Tuple[int, int]]) -> bool:
        return self.bottom - 1 > 0 and (
            [self.left, self.bottom] not in occupied
            and [self.left + 1, self.bottom - 1] not in occupied
            and [self.right, self.bottom] not in occupied
        )

    def freeze(self) -> List[int]:
        retval = [
            [x, self.bottom + 1] for x in range(self.left, self.right + 1)
        ]
        retval.extend(
            [[self.left + 1, self.top], (self.left + 1, self.bottom)]
        )
        return retval


class Ell(Shape):
    height = 3
    width = 3

    def can_blow(
        self, direction: str, occupied: List[Tuple[int, int]]
    ) -> bool:
        if direction == "<":
            return (
                self.left > 0
                and [self.left - 1, self.bottom] not in occupied
                and not any(
                    [self.right - 1, y] in occupied
                    for y in range(self.top, self.bottom, -1)
                )
            )
        return super().can_blow(direction, occupied)

    def freeze(self) -> List[int]:
        retval = [[x, self.bottom] for x in range(self.left, self.right + 1)]
        retval.extend(
            [[self.right, y] for y in range(self.top, self.bottom, -1)]
        )
        return retval


class Tall(Shape):
    height = 4
    width = 1


class Box(Shape):
    height = 2
    width = 2


def prune_occupied(occupied: List[Tuple[int, int]], max_height: int) -> None:
    to_delete = []
    for point in occupied:
        # i tried to write a proper pruning algorithm, but couldn't
        # get it to work. online i saw people pruning anything over
        # 100, which, whatever.
        if point[1] < max_height - 100:
            to_delete.append(point)
    for point in to_delete:
        occupied.remove(point)


def main():
    max_height = 0
    jets = itertools.cycle(fileinput.input()[0].strip())
    shapes = itertools.cycle([Flat, Plus, Ell, Tall, Box])
    occupied: List[Tuple[int, int]] = []

    start = time.time()
    for i in range(ROCKS):
        if i > 0 and i % 10000 == 0:
            elapsed = time.time() - start
            print(
                "%s: %0.3f years remaining"
                % (i, (elapsed / i / 60 / 60 / 24 / 365) * (ROCKS - i),)
            )
        rock = next(shapes)(max_height + 4)
        stopped = False
        while not stopped:
            rock.blow(next(jets), occupied)
            stopped, new_occupied = rock.drop(occupied)
            if stopped:
                occupied.extend(new_occupied)
                max_height = max(max_height, *[p[1] for p in new_occupied])
        prune_occupied(occupied, max_height)
    print(max_height)


if __name__ == "__main__":
    main()
