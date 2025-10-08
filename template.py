#!/usr/bin/env python

import aocd
from typing import Any

def part1(input_data: str, **extra: Any) -> aocd.types.AnswerValue:
    pass

def part2(input_data: str, **extra: Any) -> aocd.types.AnswerValue:
    pass

def main() -> None:
    day, year = aocd.get.get_day_and_year()
    puzzle = aocd.models.Puzzle(year, day)

    for example in puzzle.examples:
        assert part1(example.input_data, **example.extra) == answer_a
    puzzle.answer_a = part1(puzzle.input_data)

    for example in puzzle.examples:
        assert part2(example.input_data, **example.extra) == answer_b
    puzzle.answer_b = part2(puzzle.input_data)


if __name__ == "__main__":
    main()
