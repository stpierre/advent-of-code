#!/usr/bin/env python

import asyncio
import importlib
import logging

from loguru import logger
from temporalio import client, worker

import common
import workflows


async def main() -> None:
    logging.basicConfig(level=logging.INFO)

    all_workflows = [
        workflows.SolvePart,
        workflows.Solve,
        workflows.RunExamples,
    ]
    all_activities = [
        workflows.Solve().fetch_input_data,
        workflows.Solve().set_answer,
        workflows.RunExamples().fetch_examples,
    ]

    for day in range(1, 26):
        try:
            mod = importlib.import_module(f"solutions.day{day:02}")
        except ImportError:
            logger.info("No module found for day {}", day)
        else:
            for part in workflows.Part:
                if wf := getattr(mod, part.value.class_name, None):
                    logger.debug(
                        "Discovered workflow {}.{}",
                        mod.__name__,
                        part.value.class_name,
                    )
                    all_workflows.append(wf)

    temporal_client = await client.Client.connect(
        "localhost:7233",
        namespace="default",
    )
    logger.info("Loading worker with:")
    logger.info("  Workflows: {}", all_workflows)
    logger.info("  Activities: {}", all_activities)
    temporal_worker = worker.Worker(
        temporal_client,
        task_queue=common.TASK_QUEUE_NAME,
        workflows=all_workflows,
        activities=all_activities,
    )
    await temporal_worker.run()


if __name__ == "__main__":
    asyncio.run(main())
