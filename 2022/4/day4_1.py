#!/usr/bin/env python

import fileinput
from typing import Tuple


def parse_line(line: str) -> Tuple[int, int, int, int]:
    range1, range2 = line.strip().split(",")
    start1, end1 = range1.split("-")
    start2, end2 = range2.split("-")
    return (int(start1), int(end1), int(start2), int(end2))


def main():
    count = 0
    for line in fileinput.input():
        start1, end1, start2, end2 = parse_line(line)
        if (start1 <= start2 and end1 >= end2) or (
            start2 <= start1 and end2 >= end1
        ):
            count += 1
    print(count)


if __name__ == "__main__":
    main()
