import dataclasses
import datetime
from typing import Any

import temporalio.common


@dataclasses.dataclass
class PartInput:
    input_data: str
    extra: dict[str, Any] | None = None


NEVER_RETRY = temporalio.common.RetryPolicy(maximum_attempts=1)

YEAR = 2020
TASK_QUEUE_NAME = f"aoc-{YEAR}"

AOC_API_ACTIVITY_TIMEOUT = datetime.timedelta(5)
AOC_API_RETRY = temporalio.common.RetryPolicy(maximum_attempts=1)
