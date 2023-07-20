#!/usr/bin/env python

import fileinput


def main():
    cycle = 0
    X = 1

    for line in fileinput.input():
        args = line.strip().split()
        cmd = args[0]
        cycles = 1 if cmd == "noop" else 2

        for _ in range(cycles):
            if X - 1 <= cycle % 40 <= X + 1:
                print("#", end="")
            else:
                print(".", end="")
            if (cycle + 1) % 40 == 0:
                print()

            cycle += 1

        if cmd == "addx":
            value = int(args[1])
            X += value


if __name__ == "__main__":
    main()
