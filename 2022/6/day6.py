#!/usr/bin/env python

import collections
import fileinput


def main():
    marker_length = 14
    # if the file was bigger we'd want to read it in batches of
    # characters, but it's small enough we don't have to worry about
    # memeory
    line = fileinput.input()[0]
    chars = collections.deque(maxlen=marker_length)
    for i, char in enumerate(line):
        chars.append(char)
        if len(set(chars)) == marker_length:
            print(i + 1)
            break


if __name__ == "__main__":
    main()
