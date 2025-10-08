#!/usr/bin/env python

import asyncio
import dataclasses
import enum
import re
import sys
import time
from collections.abc import Callable
from typing import Any

import aocd
import click
import pytest
from loguru import logger

DEBUG = True


@dataclasses.dataclass
class CubeSet:
    red: int
    blue: int
    green: int

    def power(self) -> int:
        return self.red * self.blue * self.green

    def __str__(self) -> str:
        return f"{self.red} red, {self.blue} blue, {self.green} green"


class Handful(CubeSet):
    _cube_re = re.compile(r"^(?P<count>\d+) (?P<color>red|green|blue)")

    def __init__(self, cube_data: str) -> None:
        super().__init__(0, 0, 0)
        for cubes in cube_data.split(", "):
            if match := self._cube_re.match(cubes):
                setattr(self, match.group("color"), int(match.group("count")))

    async def is_possible(self, red: int, green: int, blue: int) -> bool:
        return self.red <= red and self.green <= green and self.blue <= blue


class Game:
    _game_num_re = re.compile(r"^Game (\d+)")

    def __init__(self, line: str) -> None:
        game_info, handsful = line.split(": ")
        self.game_id = int(self._game_num_re.match(game_info).group(1))
        self.handfuls = [Handful(h) for h in handsful.split("; ")]

    async def is_possible(self, red: int, green: int, blue: int) -> bool:
        return all(
            await asyncio.gather(
                *[h.is_possible(red, green, blue) for h in self.handfuls]
            )
        )

    async def min_cubes(self) -> CubeSet:
        result = CubeSet(0, 0, 0)
        for handful in self.handfuls:
            result.red = max(result.red, handful.red)
            result.blue = max(result.blue, handful.blue)
            result.green = max(result.green, handful.green)
        return result


async def part1(input_data: str, **extra: Any) -> aocd.types.AnswerValue:
    games = {}
    async with asyncio.TaskGroup() as tg:
        for line in input_data.splitlines():
            game = Game(line)
            games[game.game_id] = tg.create_task(game.is_possible(12, 13, 14))

    success_sum = 0
    for game_id, task in games.items():
        if task.result():
            logger.info("Game {} is possible", game_id)
            success_sum += game_id
    return success_sum


async def part2(input_data: str, **_: Any) -> aocd.types.AnswerValue:
    games = {}
    async with asyncio.TaskGroup() as tg:
        for line in input_data.splitlines():
            game = Game(line)
            games[game.game_id] = tg.create_task(game.min_cubes())

    power_sum = 0
    for game_id, task in games.items():
        min_cubes = task.result()
        power = min_cubes.power()
        logger.debug(
            "Game {} requires {}, power={}", game_id, min_cubes, power
        )
        power_sum += power
    return power_sum


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
    def func(self) -> Callable:
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


async def _run_examples(puzzle: aocd.models.Puzzle, part: _Part) -> None:
    for i, example in enumerate(puzzle.examples):
        answer = part.value.get_answer(example)
        if answer:
            logger.info("Checking {} against example #{}", part, i + 1)
            actual = await part.value.func(
                example.input_data,
                **(example.extra or {}),
            )
            assert actual == answer or str(actual) == answer, (
                f"Wrong answer for {part} example #{i + 1}: "
                f"{actual!r} != {example.answer_a!r}"
            )


async def _run(*, fast: bool = False) -> None:
    day, year = aocd.get.get_day_and_year()
    logger.info("Fetching puzzle data for {}/{}", year, day)
    puzzle = aocd.models.Puzzle(year, day)

    if not fast:
        await _run_examples(puzzle, _Part.ONE)
    logger.info("Running part 1")
    start = time.time()
    puzzle.answer_a = await part1(puzzle.input_data)
    logger.info("Ran part 1 in {}s", time.time() - start)

    if not fast:
        await _run_examples(puzzle, _Part.TWO)
    logger.info("Running part 2")
    start = time.time()
    puzzle.answer_b = await part2(puzzle.input_data)
    logger.info("Ran part 2 in {}s", time.time() - start)


@click.command()
@click.option(
    "-d/-q", "--debug/--quiet", help="Enable/disable debugging", default=DEBUG
)
@click.option(
    "-t",
    "--test-only",
    help="Only run tests",
    is_flag=True,
    default=False,
)
@click.option(
    "-f",
    "--fast",
    help="Only run the actual problem, skipping examples and unit tests",
    is_flag=True,
    default=False,
)
def main(debug: bool, test_only: bool, fast: bool) -> None:  # noqa: FBT001
    if not debug:
        logger.remove()
        logger.add(sys.stderr, level="INFO")

    if not fast:
        _run_tests()
        if test_only:
            sys.exit(0)
    asyncio.run(_run(fast=fast))


if __name__ == "__main__":
    main()
