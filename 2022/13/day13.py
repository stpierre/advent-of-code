#!/usr/bin/env python

import fileinput
import functools
import itertools
import json
import operator


def compare(lhs, rhs):
    if isinstance(lhs, list) and isinstance(rhs, list):
        for i in range(min(len(lhs), len(rhs))):
            result = compare(lhs[i], rhs[i])
            if result != 0:
                return result
        if len(lhs) > len(rhs):
            return 1
        if len(rhs) > len(lhs):
            return -1
        return 0
    if isinstance(lhs, list):
        return compare(lhs, [rhs])
    if isinstance(rhs, list):
        return compare([lhs], rhs)
    if lhs < rhs:
        return -1
    if lhs > rhs:
        return 1
    return 0


def main():
    packet1 = None
    packet2 = None
    ordered = []
    packets = [[[2]], [[6]]]
    pair = 1
    for line in fileinput.input():
        if not line.strip():
            if compare(packet1, packet2) == -1:
                ordered.append(pair)
            pair += 1
            packet1 = packet2 = None
        elif packet1 is None:
            packet1 = json.loads(line)
            packets.append(packet1)
        else:
            packet2 = json.loads(line)
            packets.append(packet2)
    if compare(packet1, packet2):
        ordered.append(pair)
    print("part 1: %s" % sum(ordered))

    divider_indexes = []
    for idx, packet in enumerate(
        sorted(packets, key=functools.cmp_to_key(compare))
    ):
        if packet in ([[2]], [[6]]):
            divider_indexes.append(idx + 1)
    print(
        "part 2: %s"
        % list(itertools.accumulate(divider_indexes, operator.mul))[-1]
    )


if __name__ == "__main__":
    main()
