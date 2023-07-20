#!/usr/bin/env python

import fileinput
from typing import Set, Tuple


def parse_line(line: str) -> Tuple[Set[int], Set[int]]:
    range1, range2 = line.strip().split(",")
    start1, end1 = range1.split("-")
    start2, end2 = range2.split("-")
    return (
        set(range(int(start1), int(end1) + 1)),
        set(range(int(start2), int(end2) + 1)),
    )


def main():
    count = 0
    for line in fileinput.input():
        range1, range2 = parse_line(line)
        if range1 & range2:
            count += 1
    print(count)


if __name__ == "__main__":
    main()
