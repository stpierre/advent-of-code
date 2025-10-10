from temporalio import workflow

import common


@workflow.defn
class Part1:
    @workflow.run
    async def run(self, data: common.PartInput) -> str:
        entries = [int(line) for line in data.input_data.splitlines()]
        for i, first in enumerate(entries):
            for second in entries[i + 1 :]:
                result = first + second
                if result == common.YEAR:
                    answer = first * second
                    workflow.logger.info(
                        "Found matching pair: %s + %s = %s, %s * %s = %s",
                        first,
                        second,
                        common.YEAR,
                        first,
                        second,
                        answer,
                    )
                    return str(answer)
                workflow.logger.debug(
                    "Found non-matching pair: %s + %s = %s",
                    first,
                    second,
                    first + second,
                )
        raise Exception("No matching entries found")


@workflow.defn
class Part2:
    @workflow.run
    async def run(self, data: common.PartInput) -> str:
        entries = [int(line) for line in data.input_data.splitlines()]
        for i, first in enumerate(entries):
            for j, second in enumerate(entries[i + 1 :]):
                for third in entries[j + 1 :]:
                    result = first + second + third
                    if result == common.YEAR:
                        answer = first * second * third
                        workflow.logger.info(
                            "Found matching pair: %s + %s + %s = %s, "
                            "%s * %s * %s = %s",
                            first,
                            second,
                            third,
                            common.YEAR,
                            first,
                            second,
                            third,
                            answer,
                        )
                        return str(answer)
                    workflow.logger.debug(
                        "Found non-matching pair: %s + %s + %s = %s",
                        first,
                        second,
                        third,
                        result,
                    )
        raise Exception("No matching entries found")
