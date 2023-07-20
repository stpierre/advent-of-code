#!/usr/bin/env python

import fileinput


def main():
    cycle = 0
    X = 1
    result = 0

    for line in fileinput.input():
        args = line.strip().split()
        cmd = args[0]
        cycles = 1 if cmd == "noop" else 2

        for _ in range(cycles):
            cycle += 1
            if cycle in (20, 60, 100, 140, 180, 220):
                result += cycle * X

        if cmd == "addx":
            value = int(args[1])
            X += value

        if cycle > 220:
            break
    print(result)


if __name__ == "__main__":
    main()
