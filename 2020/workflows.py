import dataclasses
import enum
import importlib
from typing import Any

from temporalio import activity, workflow

import common

with workflow.unsafe.imports_passed_through():
    import aocd


@dataclasses.dataclass(order=True, frozen=True, eq=True)
class _PartDescriptor:
    name: str
    class_name: str
    answer_name: str

    def get_answer(self, example: aocd.examples.Example) -> str:
        return getattr(example, f"answer_{self.answer_name}")

    @property
    def task_id(self) -> str:
        return self.name.replace(" ", "")

    def __str__(self) -> str:
        return self.name


class Part(enum.Enum):
    ONE = _PartDescriptor("part 1", "Part1", "a")
    TWO = _PartDescriptor("part 2", "Part2", "b")

    def __str__(self) -> str:
        return str(self.value)


def _get_workflow_class(part: _PartDescriptor, day: int) -> type:
    mod = importlib.import_module(f"solutions.day{day:02}")
    try:
        return getattr(mod, part.class_name)
    except AttributeError:
        workflow.logger.warning(
            "No solution workflow defined for day %s %s", day, part
        )
        raise


@dataclasses.dataclass
class SolveInput:
    day: int
    session_token: str
    fast: bool = False

    @property
    def task_id(self) -> str:
        suffix = "-fast" if self.fast else ""
        return f"solve-day{self.day}-{self.session_token[0:6]}{suffix}"


@dataclasses.dataclass
class SolveOutput:
    answers: list[tuple[_PartDescriptor, str | None]]


@dataclasses.dataclass
class _SetAnswerInput:
    problem: SolveInput
    part: _PartDescriptor
    answer: str


@workflow.defn
class Solve:
    @workflow.run
    async def run(self, data: SolveInput) -> SolveOutput:
        input_data = await workflow.execute_activity_method(
            self.fetch_input_data,
            data,
            retry_policy=common.AOC_API_RETRY,
            start_to_close_timeout=common.AOC_API_ACTIVITY_TIMEOUT,
        )

        answers = {}
        for part in Part:
            part_input = SolvePartInput(
                data,
                part.value,
                input_data,
            )
            answers[part.value] = await workflow.execute_child_workflow(
                SolvePart,
                part_input,
                task_queue=common.TASK_QUEUE_NAME,
                id=part_input.task_id,
            )

            if answers[part.value] is not None:
                await workflow.execute_activity_method(
                    self.set_answer,
                    _SetAnswerInput(data, part.value, answers[part.value]),
                    start_to_close_timeout=common.AOC_API_ACTIVITY_TIMEOUT,
                )
        return SolveOutput(list(answers.items()))

    @activity.defn
    async def fetch_input_data(self, data: SolveInput) -> str:
        return aocd.get_data(
            session=data.session_token, day=data.day, year=common.YEAR
        )

    @activity.defn
    async def set_answer(self, data: _SetAnswerInput) -> None:
        aocd.submit(
            answer=data.answer,
            part=data.part.answer_name,
            day=data.problem.day,
            year=common.YEAR,
            session=data.problem.session_token,
        )


@dataclasses.dataclass
class SolvePartInput:
    problem: SolveInput
    part: _PartDescriptor
    input_data: str

    @property
    def task_id(self) -> str:
        return f"{self.problem.task_id}-{self.part.task_id}"


@workflow.defn
class SolvePart:
    @workflow.run
    async def run(self, data: SolvePartInput) -> str | None:
        if not data.problem.fast:
            ex_input = RunExamplesInput(data.problem, data.part)
            await workflow.execute_child_workflow(
                RunExamples,
                ex_input,
                id=ex_input.task_id,
                retry_policy=common.NEVER_RETRY,
            )

        try:
            cls = _get_workflow_class(data.part, data.problem.day)
        except AttributeError:
            return None

        return await workflow.execute_child_workflow(
            cls,
            common.PartInput(data.input_data),
            id=f"run-{data.task_id}",
            task_queue=common.TASK_QUEUE_NAME,
        )


@dataclasses.dataclass
class RunExamplesInput:
    problem: SolveInput
    part: _PartDescriptor

    @property
    def task_id(self) -> str:
        return f"examples-{self.problem.task_id}-{self.part.task_id}"


@dataclasses.dataclass
class Example:
    """AOC example class.

    AOCD returns examples as NamedTuples, which get cast to lists when
    JSONified. This provides a way to turn them into dataclasses,
    which get turned into dicts.

    This also just retains a single answer -- the answer to the part
    we care about -- to reduce the data passed back from the
    fetch_examples activity.
    """

    input_data: str
    answer: str | None = None
    extra: dict[str, Any] | None = None


@workflow.defn
class RunExamples:
    @workflow.run
    async def run(self, data: RunExamplesInput) -> None:
        try:
            cls = _get_workflow_class(data.part, data.problem.day)
        except AttributeError:
            return

        examples = await workflow.execute_activity_method(
            self.fetch_examples,
            data,
            retry_policy=common.AOC_API_RETRY,
            start_to_close_timeout=common.AOC_API_ACTIVITY_TIMEOUT,
        )
        workflows = []
        for i, example in enumerate(examples):
            if example.answer is not None:
                workflow.logger.info(
                    "Checking %s against example #%s", data.part, i + 1
                )
                workflow.logger.debug(
                    "Example #%s input: %r", i + 1, example.input_data
                )
                workflow.logger.debug(
                    "Example #%s expected: %r", i + 1, example.answer
                )
                ex_data = common.PartInput(example.input_data, example.extra)
                workflow.logger.debug(
                    "Invoking child workflow %s with %s",
                    cls,
                    ex_data,
                )
                workflows.append(
                    (
                        workflow.execute_child_workflow(
                            cls,
                            ex_data,
                            task_queue=common.TASK_QUEUE_NAME,
                            id=f"example-{i}-{data.task_id}",
                        ),
                        example.answer,
                    ),
                )
        for wf, expected in workflows:
            actual = await wf
            assert actual == expected, (
                f"Wrong answer for {data.part} example #{i + 1}: "
                f"{actual!r} != {expected!r}"
            )

    @activity.defn
    async def fetch_examples(self, data: RunExamplesInput) -> list[Example]:
        retval = [
            Example(ex.input_data, data.part.get_answer(ex), ex.extra)
            for ex in aocd.get_puzzle(
                session=data.problem.session_token,
                day=data.problem.day,
                year=common.YEAR,
            ).examples
        ]
        activity.logger.info("Fetched %s examples: %s", len(retval), retval)
        return retval
