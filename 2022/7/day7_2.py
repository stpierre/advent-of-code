#!/usr/bin/env python

import fileinput

import anytree


class DEntry(anytree.NodeMixin):
    def __init__(self, name: str, size: int = 0, parent=None, children=None):
        self.name = name
        self.size = size
        self.parent = parent
        if children:
            self.children = children


def get_sizes(node, sizes):
    if node.size:
        return node.size

    size = sum(get_sizes(child, sizes) for child in node.children)
    sizes[node] = size
    return size


def main():
    root = None
    node = None

    for line in fileinput.input():
        parts = line.strip().split()
        if parts[0:2] == ["$", "cd"]:
            if parts[2] == "..":
                node = node.parent
            elif parts[2] == "/":
                root = node = DEntry(parts[2])
            else:
                node = next(n for n in node.children if n.name == parts[2])
        elif parts[0:2] == ["$", "ls"]:
            pass
        elif parts[0] == "dir":
            DEntry(parts[1], parent=node)
        elif parts[0].isdecimal():
            DEntry(parts[1], size=int(parts[0]), parent=node)

    # do we use a nasty side-effect of passing dicts by reference? you
    # freaking bet we do. immutable programming is for the *birds*
    sizes = {}
    total = get_sizes(root, sizes)

    needed = 30000000 - (70000000 - total)

    to_delete = None
    for dir_size in sizes.values():
        if dir_size > needed:
            if to_delete is None:
                to_delete = dir_size
            else:
                to_delete = min(to_delete, dir_size)
    print(to_delete)


if __name__ == "__main__":
    main()
