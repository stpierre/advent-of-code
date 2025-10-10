#!/usr/bin/env python

import asyncio

import click
from loguru import logger
from temporalio import client

import common
import workflows


async def _main(day: int, session_token: str, *, fast: bool = False) -> None:
    temporal_client = await client.Client.connect("localhost:7233")
    data = workflows.SolveInput(day, session_token, fast)
    logger.debug("Running {} with {}", workflows.Solve, data)
    try:
        output = await temporal_client.execute_workflow(
            workflows.Solve,
            data,
            id=data.task_id,
            task_queue=common.TASK_QUEUE_NAME,
        )
    except client.WorkflowFailureError as err:
        logger.error("Workflow failed: {}", err)
    else:
        for part, answer in output.answers:
            logger.info("Day {} {}: {}", day, part, answer)


@click.command()
@click.option("--aoc-session", envvar="AOC_SESSION", help="AOC session secret")
@click.option(
    "-f",
    "--fast",
    help="Only run the actual problem, skipping examples and unit tests",
    is_flag=True,
    default=False,
)
@click.argument("day", type=int)
def cli(aoc_session: str, fast: bool, day: int) -> None:
    asyncio.run(_main(day, aoc_session, fast=fast))


if __name__ == "__main__":
    cli()
