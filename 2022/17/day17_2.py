#!/usr/bin/env python

import collections
import fileinput
import itertools
from typing import List, Tuple

WIDTH = 7
ROCKS = 1000000000000


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
        if direction == ">":
            return self.right < WIDTH - 1 and not any(
                (self.right + 1, y) in occupied
                for y in range(self.bottom, self.top + 1)
            )
        if direction == "<":
            return self.left > 0 and not any(
                (self.left - 1, y) in occupied
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
            (x, self.bottom - 1) in occupied
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
                retval.append((x, y))
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
            return (
                self.right < WIDTH - 1
                and (self.right + 1, self.bottom + 1) not in occupied
                and (self.right, self.bottom) not in occupied
                and (self.right, self.top) not in occupied
            )

        if direction == "<":
            return (
                self.left > 0
                and (self.left - 1, self.bottom + 1) not in occupied
                and (self.left, self.bottom) not in occupied
                and (self.left, self.top) not in occupied
            )
        return super().can_blow(direction, occupied)

    def can_drop(self, occupied: List[Tuple[int, int]]) -> bool:
        return self.bottom - 1 > 0 and (
            (self.left, self.bottom) not in occupied
            and (self.left + 1, self.bottom - 1) not in occupied
            and (self.right, self.bottom) not in occupied
        )

    def freeze(self) -> List[int]:
        retval = [
            (x, self.bottom + 1) for x in range(self.left, self.right + 1)
        ]
        retval.extend(
            [(self.left + 1, self.top), (self.left + 1, self.bottom)]
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
                and (self.left - 1, self.bottom) not in occupied
                and not any(
                    (self.right - 1, y) in occupied
                    for y in range(self.top, self.bottom, -1)
                )
            )
        return super().can_blow(direction, occupied)

    def freeze(self) -> List[int]:
        retval = [(x, self.bottom) for x in range(self.left, self.right + 1)]
        retval.extend(
            [(self.right, y) for y in range(self.top, self.bottom, -1)]
        )
        return retval


class Tall(Shape):
    height = 4
    width = 1


class Box(Shape):
    height = 2
    width = 2


def prune_occupied(
    occupied: List[Tuple[int, int]], max_height: int
) -> List[Tuple[int, int]]:
    new_occupied = []
    for point in occupied:
        # i tried to write a proper pruning algorithm, but couldn't
        # get it to work. online i saw people pruning anything over
        # 100, which, whatever.
        if point[1] > max_height - 100:
            new_occupied.append(point)
    return new_occupied


def main():
    max_height = 0
    jets = itertools.cycle(fileinput.input()[0].strip())
    shapes = itertools.cycle([Flat, Plus, Ell, Tall, Box])
    occupied: List[Tuple[int, int]] = []
    heights: List[int] = [0 for _ in range(WIDTH)]
    cycle_starts = collections.defaultdict(list)
    found_cycle = False

    rock_num = 0
    while rock_num < ROCKS:
        if rock_num > 0 and rock_num % 100 == 0:
            print(rock_num)
        rock = next(shapes)(max_height + 4)
        stopped = False
        while not stopped:
            rock.blow(next(jets), occupied)
            stopped, new_occupied = rock.drop(occupied)
            if stopped:
                occupied.extend(new_occupied)
                occupied = prune_occupied(occupied, max_height)
                max_height = max(max_height, *[p[1] for p in new_occupied])
                if not found_cycle:
                    for point in new_occupied:
                        heights[point[0]] = max(heights[point[0]], point[1])

                    height_deltas = [
                        max_height - heights[x] for x in range(WIDTH)
                    ]
                    if all(-1 <= d <= 1 for d in height_deltas):
                        key = (type(rock).__name__, tuple(height_deltas))
                        print(
                            "found potential cycle at %s: %s" % (rock_num, key)
                        )
                        cycle_starts[key].append((rock_num, max_height))
                        if len(cycle_starts[key]) >= 3:
                            print(
                                "  found key %s times (%s), investigating"
                                % (
                                    len(cycle_starts[key]),
                                    [c[0] for c in cycle_starts[key]],
                                )
                            )
                            cycle_length = None
                            for i in range(len(cycle_starts[key]) - 1):
                                cur_dist = (
                                    cycle_starts[key][i + 1][0]
                                    - cycle_starts[key][i][0]
                                )
                                if cycle_length is None:
                                    cycle_length = cur_dist
                                elif cur_dist != cycle_length:
                                    break
                            else:
                                height_added = (
                                    cycle_starts[key][1][1]
                                    - cycle_starts[key][0][1]
                                )
                                cycles_needed = (
                                    ROCKS - rock_num
                                ) // cycle_length
                                print(
                                    f"  found cycle of length {cycle_length}, "
                                    f"adding {height_added} to height. "
                                    f"{cycles_needed} cycles needed"
                                )
                                found_cycle = True
                                max_height += height_added * cycles_needed
                                rock_num += cycle_length * cycles_needed
                                occupied = [
                                    (max_height - d, y)
                                    for y, d in enumerate(height_deltas)
                                ]
                                print(
                                    "  fast-forwarding to rock #%s, max height=%s"
                                    % (rock_num, max_height)
                                )
                                print(occupied)
        rock_num += 1
    print(max_height)


if __name__ == "__main__":
    main()
