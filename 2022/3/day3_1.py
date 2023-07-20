#!/usr/bin/env python

import fileinput


def main():
    total = 0
    for line in fileinput.input():
        line = line.strip()
        compartment_size = len(line) // 2
        oops = list(
            set(line[0:compartment_size]).intersection(line[compartment_size:])
        )[0]
        if oops < "a":
            total += ord(oops) - 38
        else:
            total += ord(oops) - 96
    print(total)


if __name__ == "__main__":
    main()
