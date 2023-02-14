from typing import Tuple
from matplotlib import pyplot as plt
import argparse
import os
import os.path


def plot(data: list[list[float]], output: str):
    assert len(data) > 0, 'empty data'

    plt.scatter(list(range(len(data))), [i[0] for i in data])
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

    data: list[list[float]] = []
    with open(args.file, 'r') as file:
        lines = file.readlines()
        for line in lines:
            split = line.split(',')
            if len(split) != 3:
                assert 'invalid data'

            index = int(split[0])
            start_time = int(float(split[1]) * 1000)
            elapsed = float(split[2])

            if index >= len(data):
                data.extend([[] for _ in range(index - len(data) + 1)])

            data[index].append(elapsed)

    plot(x, data, args.output)
