#!/usr/bin/env python

import fileinput


def main():
    decoder = {
        "A": "rock",
        "B": "paper",
        "C": "scissors",
        "X": "rock",
        "Y": "paper",
        "Z": "scissors",
    }
    throw_scores = {"rock": 1, "paper": 2, "scissors": 3}

    score = 0
    for line in fileinput.input():
        opp_code, my_code = line.strip().split()
        opp_throw = decoder[opp_code]
        my_throw = decoder[my_code]

        score += throw_scores[my_throw]

        if opp_throw == my_throw:
            score += 3
        elif (
            (opp_throw == "rock" and my_throw == "paper")
            or (opp_throw == "paper" and my_throw == "scissors")
            or (opp_throw == "scissors" and my_throw == "rock")
        ):
            score += 6

    print(score)


if __name__ == "__main__":
    main()
