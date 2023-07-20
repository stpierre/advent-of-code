#!/usr/bin/env python

import operator


def main():
    ranges = sorted(
        [(10, 17), (-2, 4), (12, 44), (2, 7), (8, 8)],
        key=operator.itemgetter(0),
    )
    x_min = 0
    x_max = 20

    if ranges[0][0] > x_min:
        print("%s is not excluded" % x_min)
    exclude_end = 0
    for x_range in ranges:
        if x_range[0] > exclude_end + 1:
            print("%s is not excluded" % (exclude_end + 1))
            break
        else:
            exclude_end = x_range[1]
    else:
        if exclude_end < x_max:
            print("%s is not excluded" % x_max)


if __name__ == "__main__":
    main()
