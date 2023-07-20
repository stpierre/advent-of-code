#!/usr/bin/env python

import fileinput


def main():
    rows = []
    for line in fileinput.input():
        rows.append(line.strip())

    grid_size = len(rows)

    # count all the edges
    visible = grid_size * 4 - 4
    for x in range(1, grid_size - 1):
        for y in range(1, grid_size - 1):
            cur_tree = int(rows[y][x])

            if (
                all(
                    int(rows[y][candidate_x]) < cur_tree
                    for candidate_x in range(x)
                )
                or all(
                    int(rows[y][candidate_x]) < cur_tree
                    for candidate_x in range(x + 1, grid_size)
                )
                or all(
                    int(rows[candidate_y][x]) < cur_tree
                    for candidate_y in range(y)
                )
                or all(
                    int(rows[candidate_y][x]) < cur_tree
                    for candidate_y in range(y + 1, grid_size)
                )
            ):
                visible += 1
                continue
    print(visible)


if __name__ == "__main__":
    main()
