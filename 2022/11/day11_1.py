#!/usr/bin/env python

import collections
import fileinput
import re
from typing import List, Tuple


class Monkey:
    def __init__(
        self,
        items: List[int],
        operation: str,
        divisible: int,
        throw_to_true: int,
        throw_to_false: int,
    ) -> None:
        self.items = items
        self.operation = operation
        self.divisible = divisible
        self.throw_to_true = throw_to_true
        self.throw_to_false = throw_to_false

    def throw_one(self) -> Tuple[int, int]:
        item = eval(self.operation, {"old": self.items.pop(0)}) // 3
        if item % self.divisible == 0:
            return self.throw_to_true, item
        return self.throw_to_false, item


def parse_input() -> List[Monkey]:
    monkeys = {}

    monkey_num = None
    items = []
    operation = None
    divisible = None
    throw_to_true = None
    throw_to_false = None
    for line in fileinput.input():
        if match := re.match(r"Monkey (?P<num>\d+):", line):
            monkey_num = int(match.group("num"))
            items = []
            operation = None
            divisible = None
            throw_to_true = None
            throw_to_false = None
        elif match := re.match(
            r"\s*Starting items: (?P<items>[0-9, ]*)", line
        ):
            items = [int(i) for i in match.group("items").split(", ")]
        elif match := re.match(
            r"\s*Operation: new = (?P<expression>.*)", line.strip()
        ):
            operation = match.group("expression")
        elif line.startswith("  Test:"):
            _, raw_divisor = line.strip().rsplit(maxsplit=1)
            divisible = int(raw_divisor)
        elif line.startswith("    If true:"):
            _, target = line.strip().rsplit(maxsplit=1)
            throw_to_true = int(target)
        elif line.startswith("    If false:"):
            _, target = line.strip().rsplit(maxsplit=1)
            throw_to_false = int(target)
        elif not line.strip():
            monkeys[monkey_num] = Monkey(
                items, operation, divisible, throw_to_true, throw_to_false
            )
    monkeys[monkey_num] = Monkey(
        items, operation, divisible, throw_to_true, throw_to_false
    )
    return monkeys


def main():
    monkeys = parse_input()
    inspected = collections.defaultdict(int)
    for _ in range(20):
        for monkey_num in sorted(monkeys.keys()):
            monkey = monkeys[monkey_num]
            while monkey.items:
                inspected[monkey_num] += 1
                target, item = monkey.throw_one()
                monkeys[target].items.append(item)
    inspections = sorted(inspected.values())[-2:]
    print(inspections[0] * inspections[1])


if __name__ == "__main__":
    main()
