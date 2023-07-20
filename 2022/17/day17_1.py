#!/usr/bin/env python

import dataclasses
import fileinput
import itertools
from typing import List, Tuple

WIDTH = 7
ROCKS = 2022
DEBUG = False


@dataclasses.dataclass(frozen=True, eq=True)
class Point:
    x: int
    y: int


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

    def can_blow(self, direction: str, occupied: List[Point]) -> bool:
        if direction == ">":
            return self.right < WIDTH - 1 and not any(
                Point(self.right + 1, y) in occupied
                for y in range(self.bottom, self.top + 1)
            )
        if direction == "<":
            return self.left > 0 and not any(
                Point(self.left - 1, y) in occupied
                for y in range(self.bottom, self.top + 1)
            )

        raise Exception(f"Unknown direction {direction}")

    def blow(self, direction: str, occupied: List[Point]) -> None:
        debug("blowing rock %s" % direction)
        if not self.can_blow(direction, occupied):
            return
        if direction == ">":
            self.left += 1
        elif direction == "<":
            self.left -= 1
        else:
            raise Exception(f"Unknown direction {direction}")

    def can_drop(self, occupied: List[Point]) -> bool:
        return self.bottom - 1 > 0 and not any(
            Point(x, self.bottom - 1) in occupied
            for x in range(self.left, self.right + 1)
        )

    def drop(self, heights: List[int]) -> Tuple[bool, List[int]]:
        debug("dropping rock v")
        if not self.can_drop(heights):
            return True, self.freeze()
        self.bottom -= 1
        return False, heights

    def freeze(self) -> List[Point]:
        retval = []
        for x in range(self.left, self.right + 1):
            for y in range(self.bottom, self.top + 1):
                retval.append(Point(x, y))
        return retval


class Flat(Shape):
    height = 1
    width = 4


class Plus(Shape):
    height = 3
    width = 3

    def can_blow(self, direction: str, occupied: List[Point]) -> bool:
        if direction == ">":
            return (
                self.right < WIDTH - 1
                and Point(self.right + 1, self.bottom + 1) not in occupied
                and Point(self.right, self.bottom) not in occupied
                and Point(self.right, self.top) not in occupied
            )

        if direction == "<":
            return (
                self.left > 0
                and Point(self.left - 1, self.bottom + 1) not in occupied
                and Point(self.left, self.bottom) not in occupied
                and Point(self.left, self.top) not in occupied
            )
        return super().can_blow(direction, occupied)

    def can_drop(self, occupied: List[Point]) -> bool:
        return self.bottom - 1 > 0 and (
            Point(self.left, self.bottom) not in occupied
            and Point(self.left + 1, self.bottom - 1) not in occupied
            and Point(self.right, self.bottom) not in occupied
        )

    def freeze(self) -> List[int]:
        retval = [
            Point(x, self.bottom + 1) for x in range(self.left, self.right + 1)
        ]
        retval.extend(
            [Point(self.left + 1, self.top), Point(self.left + 1, self.bottom)]
        )
        return retval


class Ell(Shape):
    height = 3
    width = 3

    def can_blow(self, direction: str, occupied: List[Point]) -> bool:
        if direction == "<":
            return (
                self.left > 0
                and Point(self.left - 1, self.bottom) not in occupied
                and not any(
                    Point(self.right - 1, y) in occupied
                    for y in range(self.top, self.bottom, -1)
                )
            )
        return super().can_blow(direction, occupied)

    def freeze(self) -> List[int]:
        retval = [
            Point(x, self.bottom) for x in range(self.left, self.right + 1)
        ]
        retval.extend(
            [Point(self.right, y) for y in range(self.top, self.bottom, -1)]
        )
        return retval


class Tall(Shape):
    height = 4
    width = 1


class Box(Shape):
    height = 2
    width = 2


def debug(msg: str):
    if DEBUG:
        print(msg)


def dump_map(occupied, rock=None, max_height=None, force=False):
    if DEBUG or force:
        rock_points = rock.freeze() if rock else []
        top = rock.top if rock else max_height
        for y in range(top, 0, -1):
            print("|", end="")
            for x in range(0, WIDTH):
                cur_point = Point(x, y)
                if cur_point in rock_points:
                    print("@", end="")
                elif cur_point in occupied:
                    print("#", end="")
                else:
                    print(" ", end="")
            print("|")
        print("+" + "-" * WIDTH + "+")


def prune_occupied(occupied: List[Point], max_height: int) -> List[Point]:
    new_occupied = []
    for point in occupied:
        # i tried to write a proper pruning algorithm, but couldn't
        # get it to work. online i saw people pruning anything over
        # 100, which, whatever.
        if point.y > max_height - 100:
            new_occupied.append(point)
    return new_occupied


def main():
    max_height = 0
    jets = itertools.cycle(fileinput.input()[0].strip())
    shapes = itertools.cycle([Flat, Plus, Ell, Tall, Box])
    expected_heights = open("heights.txt").read().splitlines()
    occupied: List[Point] = []

    for i in range(ROCKS):
        rock = next(shapes)(max_height + 4)
        stopped = False
        dump_map(occupied, rock=rock)
        while not stopped:
            debug(
                "Rock is at %s, %s to %s, %s"
                % (rock.top, rock.left, rock.bottom, rock.right)
            )
            rock.blow(next(jets), occupied)
            dump_map(occupied, rock=rock)
            stopped, new_occupied = rock.drop(occupied)
            if stopped:
                occupied.extend(new_occupied)
                max_height = max(max_height, *[p.y for p in new_occupied])
                occupied = prune_occupied(occupied, max_height)
                dump_map(occupied, max_height=max_height)
            else:
                dump_map(occupied, rock=rock)
        if DEBUG:
            expected = int(expected_heights[i])
            if expected != max_height:
                print(
                    "Wrong height %s at %s; expected %s"
                    % (max_height, i, expected)
                )
                dump_map(occupied, max_height=max_height, force=True)
                break
    print(max_height)


if __name__ == "__main__":
    main()
