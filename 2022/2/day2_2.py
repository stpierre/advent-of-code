#!/usr/bin/env python

import fileinput

_THROW_SCORES = {"rock": 1, "paper": 2, "scissors": 3}
_WINS = {"rock": "paper", "paper": "scissors", "scissors": "rock"}
_LOSSES = {v: k for k, v in _WINS.items()}


def get_round_score(opp_throw, my_throw):
    score = _THROW_SCORES[my_throw]

    if opp_throw == my_throw:
        score += 3
    elif _WINS[opp_throw] == my_throw:
        score += 6

    return score


def main():
    throw_decoder = {
        "A": "rock",
        "B": "paper",
        "C": "scissors",
    }

    score = 0
    for line in fileinput.input():
        opp_code, result = line.strip().split()
        opp_throw = throw_decoder[opp_code]

        if result == "X":
            my_throw = _LOSSES[opp_throw]
        elif result == "Y":
            my_throw = opp_throw
        else:
            my_throw = _WINS[opp_throw]

        score += get_round_score(opp_throw, my_throw)

    print(score)


if __name__ == "__main__":
    main()
