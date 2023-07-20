#!/usr/bin/env python

import fileinput
import itertools


def main():
    draw_order = None
    boards = []

    for line in fileinput.input():
        if draw_order is None:
            draw_order = itertools.cycle(",".split(line.strip()))
        elif not line.strip():
            boards.append([])
        else:
            boards[-1].append([int(i) for i in line.strip().split()])

    print(boards)

if __name__ == "__main__":
    main()
