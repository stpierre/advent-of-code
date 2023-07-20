#!/usr/bin/env python

import fileinput
import re

import networkx


def main():
    line_re = re.compile(
        r"Valve (?P<name>[A-Z]+) has flow rate=(?P<flow_rate>\d+); tunnels lead to valves (?P<tunnels>[A-Z, ]*)"
    )
    graph = networkx.Graph()
    start = None
    edges = []

    for line in fileinput.input():
        match = line_re.match(line)
        name = match.group("name")
        graph.add_node(name, flow_rate=int(match.group("flow_rate")))
        edges.extend([(name, v) for v in match.group("tunnels").split(", ")])
        if start is None:
            start = name
    graph.add_edges_from(edges)


if __name__ == "__main__":
    main()
