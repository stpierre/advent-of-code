#!/usr/bin/env python

import aocd
import click
import dataclasses
import enum
from loguru import logger
import pytest
import sys
import time
from typing import Any

DEBUG = True


def _find_digit(line: str) -> str:
    for char in line:
        if char.isdigit():
            return char
    raise Exception(f"No digit found in {line}")


def _get_digits(line: str) -> int:
    first = _find_digit(line)
    last = _find_digit(reversed(line))
    val = f"{first}{last}"
    logger.debug("From {!r} got {}", line, val)
    return int(val)


def part1(input_data: str, **extra: Any) -> aocd.types.AnswerValue:
    return str(sum([_get_digits(line) for line in input_data.splitlines()]))


_NUMBERS = {
    "one": 1,
    "two": 2,
    "three": 3,
    "four": 4,
    "five": 5,
    "six": 6,
    "seven": 7,
    "eight": 8,
    "nine": 9,
}


def _find_number(line: str) -> str:
    for idx, char in enumerate(line):
        if char.isdigit():
            return char
        else:
            for name, num in _NUMBERS.items():
                if char == name[0] and line.find(name) == idx:
                    return num


def _rfind_number(line: str) -> str:
    for idx in range(len(line) - 1, -1, -1):
        char = line[idx]
        if char.isdigit():
            return char
        else:
            for name, num in _NUMBERS.items():
                if (idx - len(name) + 1) >= 0 and char == name[-1] and line.rfind(name) == idx - len(name) + 1:
                    return num


def _get_digits2(line: str) -> int:
    first = _find_number(line)
    last = _rfind_number(line)
    val = f"{first}{last}"
    logger.debug("From {!r} got {}", line, val)
    return int(val)


def part2(input_data: str, **extra: Any) -> aocd.types.AnswerValue:
    return str(sum([_get_digits2(line) for line in input_data.splitlines()]))


##### BOILERPLATE BELOW THIS #####

def _run_tests() -> None:
    """Run unit tests.

    This has to be done before the event loop is started, since pytest
    asyncio plugins like to start their own event loops.
    """
    logger.info("Running unit tests")
    pytest.main(["-ra", "-vv", __file__])


@dataclasses.dataclass
class _PartDescriptor:
    name: str
    func_name: str
    answer_name: str

    @property
    def func(self) -> callable:
        return globals()[self.func_name]

    def get_answer(self, example: aocd.examples.Example) -> str:
        return getattr(example, self.answer_name)

    def __str__(self) -> str:
        return self.name


class _Part(enum.Enum):
    ONE = _PartDescriptor("part 1", "part1", "answer_a")
    TWO = _PartDescriptor("part 2", "part2", "answer_b")

    def __str__(self) -> str:
        return str(self.value)


def _run_examples(puzzle: aocd.models.Puzzle, part: _Part):
    for i, example in enumerate(puzzle.examples):
        answer = part.value.get_answer(example)
        if answer:
            logger.info("Checking {} against example #{}", part, i + 1)
            actual = part.value.func(example.input_data, **(example.extra or {}))
            assert actual == answer, (
                f"Wrong answer for {part} example #{i + 1}: {actual!r} != {example.answer_a!r}"
            )


def _run(fast: bool = False) -> None:
    day, year = aocd.get.get_day_and_year()
    logger.info("Fetching puzzle data for {}/{}", year, day)
    puzzle = aocd.models.Puzzle(year, day)

    if not fast:
        _run_examples(puzzle, _Part.ONE)
    logger.info("Running part 1")
    start = time.time()
    puzzle.answer_a = part1(puzzle.input_data)
    logger.info("Ran part 1 in {}s", time.time() - start)

    if not fast:
        _run_examples(puzzle, _Part.TWO)
    logger.info("Running part 2")
    start = time.time()
    puzzle.answer_b = part2(puzzle.input_data)
    logger.info("Ran part 2 in {}s", time.time() - start)


@click.command()
@click.option("-d/-q", "--debug/--quiet", help="Enable/disable debugging", default=DEBUG)
@click.option("-t", "--test-only", help="Only run tests", is_flag=True, default=False)
@click.option("-f", "--fast", help="Only run the actual problem, skipping examples and unit tests", is_flag=True, default=False)
def main(debug: bool, test_only: bool, fast: bool) -> None:
    if not debug:
        logger.remove()
        logger.add(sys.stderr, level="INFO")

    if not fast:
        _run_tests()
        if test_only:
            sys.exit(0)
    _run(fast)

if __name__ == "__main__":
    main()
