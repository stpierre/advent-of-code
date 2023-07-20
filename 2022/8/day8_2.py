#!/usr/bin/env python

import fileinput


def main():
    rows = []
    for line in fileinput.input():
        rows.append(line.strip())

    grid_size = len(rows)

    max_scenic = 0
    # we can skip the edges since they're guaranteed to have a 0 in at
    # least one direction
    for x in range(1, grid_size - 1):
        for y in range(1, grid_size - 1):
            cur_tree = int(rows[y][x])

            left = 0
            for view_x in range(x - 1, -1, -1):
                left += 1
                if int(rows[y][view_x]) >= cur_tree:
                    break

            right = 0
            for view_x in range(x + 1, grid_size):
                right += 1
                if int(rows[y][view_x]) >= cur_tree:
                    break

            top = 0
            for view_y in range(y - 1, -1, -1):
                top += 1
                if int(rows[view_y][x]) >= cur_tree:
                    break

            bottom = 0
            for view_y in range(y + 1, grid_size):
                bottom += 1
                if int(rows[view_y][x]) >= cur_tree:
                    break

            max_scenic = max(max_scenic, left * right * top * bottom)

    print(max_scenic)


if __name__ == "__main__":
    main()
