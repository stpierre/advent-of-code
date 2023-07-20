#!/usr/bin/env python

import collections
import fileinput


def main():
    cavern = collections.defaultdict(dict)
    y_max = 0
    for line in fileinput.input():
        formation_points = line.strip().split(" -> ")
        for i, raw_coords1 in enumerate(formation_points[:-1]):
            coords1 = raw_coords1.split(",")
            x1 = int(coords1[0])
            y1 = int(coords1[1])

            coords2 = formation_points[i + 1].split(",")
            x2 = int(coords2[0])
            y2 = int(coords2[1])

            y_max = max(y_max, y1, y2)

            x_dir = 1 if x1 <= x2 else -1
            y_dir = 1 if y1 <= y2 else -1
            for x in range(x1, x2 + x_dir, x_dir):
                for y in range(y1, y2 + y_dir, y_dir):
                    cavern[x][y] = True

    sand_count = 0
    done = False
    while not done:
        sand_x = 500
        sand_y = 0
        while True:
            if sand_y >= y_max:
                done = True
                break

            if not cavern[sand_x].get(sand_y + 1):
                sand_y += 1
            elif not cavern[sand_x - 1].get(sand_y + 1):
                sand_x -= 1
                sand_y += 1
            elif not cavern[sand_x + 1].get(sand_y + 1):
                sand_x += 1
                sand_y += 1
            else:
                cavern[sand_x][sand_y] = True
                sand_count += 1
                break

    print(sand_count)


if __name__ == "__main__":
    main()
