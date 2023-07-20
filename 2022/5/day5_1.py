#!/usr/bin/env python

import collections
import fileinput
import operator
import re


def main():
    moves_started = False
    stacks = collections.defaultdict(list)
    move_re = re.compile(
        r"^move (?P<count>\d+) from (?P<src>\d) to (?P<dest>\d)"
    )

    for line in fileinput.input():
        if not line.strip():
            moves_started = True
            continue

        if not moves_started:
            # we rely on a couple of quirks of the input here, namely:
            # 1. there are only 9 stacks and all the crates use
            #    single-letter indicators, so the spacing between them
            #    stays constant: every crate except the last one takes
            #    four characters, and the last crate on a line takes
            #    three
            # 2. secondly, every line is padded with spaces to the
            #    full line length, so we can read any line and know
            #    how many stacks there are.
            # 3. the line with the stack numbers is superfluous, since
            #    the stacks start with 1 and increase monotonically
            if line.startswith(" 1"):
                continue
            for i in range(0, len(line), 4):
                crate = line[i : i + 3].strip().lstrip("[").rstrip("]")
                if crate:
                    stacks[str(i // 4 + 1)].insert(0, crate)
        else:
            match = move_re.match(line)
            for _ in range(int(match.group("count"))):
                stacks[match.group("dest")].append(
                    stacks[match.group("src")].pop()
                )

    print(
        "".join(
            [
                stack[-1]
                for _, stack in sorted(
                    stacks.items(), key=operator.itemgetter(0)
                )
            ]
        )
    )


if __name__ == "__main__":
    main()
