from typing import Tuple
from matplotlib import pyplot as plt
import argparse
import os
import os.path


def plot(x: list[int], y: list[list[float]], output: str):
    assert len(x) > 0, 'empty data'
    assert len(x) == len(y), 'invalid data'

    plt.scatter(list(range(len(x))), [i[0] for i in y])
    print(x)
    print(y)
    plt.savefig(output)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(prog='plot', description='plot tool')
    parser.add_argument('file', type=str, help='timer log file')
    parser.add_argument('output', type=str,
                        help='output image path')
    args = parser.parse_args()

    dir = os.path.dirname(args.output)
    if dir != '':
        os.makedirs(dir, exist_ok=True)

    x: list[int] = []
    y: list[list[float]] = []
    temp_y: list[float] = []
    with open(args.file, 'r') as file:
        lines = file.readlines()
        for line in lines:
            split = line.split(',')
            if len(split) != 2:
                assert 'invalid data'

            start_time = int(float(split[0]) * 1000)
            elapsed = float(split[1])

            if len(x) == 0:
                x.append(start_time)
                temp_y.append(elapsed)
            else:
                if x[-1] == start_time:
                    temp_y.append(elapsed)
                else:
                    y.append(temp_y)
                    temp_y.clear()
                    x.append(start_time)
                    temp_y.append(elapsed)
        y.append(temp_y)

    plot(x, y, args.output)
