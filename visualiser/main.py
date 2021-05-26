#!/usr/bin/python3

from argparse import ArgumentParser
from os.path import abspath, isabs, isdir, join
from os import walk
from typing import Dict, List, Tuple
from re import compile as re_compile
from math import ceil
from matplotlib import pyplot as plt
import seaborn as sns


def visualise(data: Dict[Tuple[int, int], int], sink: str, title: str) -> bool:
    if not data:
        return False

    try:
        tmp = [(f'{start}ms - {end}ms', count)
               for (start, end), count in data.items()]
        x = [t[0] for t in tmp]
        y = [t[1] for t in tmp]

        with plt.style.context('seaborn-darkgrid'):
            fig = plt.Figure(
                figsize=(16, 9),
                dpi=100)

            sns.barplot(
                x=x,
                y=y,
                orient='v',
                ax=fig.gca())

            for j, k in enumerate(fig.gca().patches):
                fig.gca().text(k.get_x() + k.get_width() / 2,
                               k.get_y() + k.get_height() * .2 + .1,
                               y[j],
                               ha='center',
                               rotation=0,
                               fontsize=12,
                               color='black')

            fig.gca().set_xlabel('Message Received After Delay',
                                 labelpad=12)
            fig.gca().set_ylabel('#-of Messages',
                                 labelpad=12)
            fig.gca().set_title(title,
                                pad=16,
                                fontsize=20)

            fig.savefig(
                sink,
                bbox_inches='tight',
                pad_inches=.5)
            plt.close(fig)

        return True
    except Exception as e:
        print(f'Error : ${e}')
        return False


def aggregated_count_by_slot(slots: List[Tuple[int, int]], bucket: Dict[int, int]) -> Dict[Tuple[int, int], int]:
    splitted = {}
    for slot in slots:
        total = 0
        (start, end) = slot
        while start <= end:
            if start in bucket:
                total += bucket[start]

            start += 1

        splitted[slot] = total

    return splitted


def splitted_delay_spectrum(delays: List[int], slot: int) -> List[Tuple[int, int]]:
    min_delay = min(delays)
    max_delay = max(delays)
    skip_by = ceil((max_delay - min_delay) / slot)

    slots = []
    while len(slots) < slot:
        next_delay = min_delay + skip_by
        if next_delay > max_delay:
            slots.append((min_delay, max_delay))
            break

        slots.append((min_delay, next_delay))
        min_delay += (skip_by + 1)

    return slots


def accumulate_data(file: str, bucket: Dict[int, int]):
    with open(file) as fd:
        while True:
            ln = fd.readline()
            if not ln:
                break

            ts = [i.strip() for i in ln.split(';')][:2]
            sent = int(ts[0])
            received = int(ts[1])
            diff = received - sent
            if diff not in bucket:
                bucket[diff] = 1
            else:
                bucket[diff] = bucket[diff] + 1

    return


def find_files(dir: str) -> List[str]:
    found = []
    reg = re_compile(r'^(log_\d+\.csv)$')
    for (dirpath, _, files) in walk(dir):
        if dirpath == dir:
            for file in files:
                if reg.match(file):
                    found.append(join(dirpath, file))

            break

    return found


def main():
    parser = ArgumentParser(description='Visualise `pub0sub` performance data')
    parser.add_argument('path', type=str, help='path to data directory')
    parser.add_argument(
        'slot', type=int, help='split delay spectrum into slots')
    args = parser.parse_args()

    target_dir = None
    if not isabs(args.path):
        target_dir = abspath(args.path)

    if not isdir(target_dir):
        print('Expected path to directory !')
        return

    found = find_files(target_dir)
    if not found:
        print('Directory walk found no file !')
        return

    bucket = {}
    for file in found:
        accumulate_data(file, bucket)

    slots = splitted_delay_spectrum(bucket.keys(), args.slot)
    print(visualise(aggregated_count_by_slot(slots, bucket), 'out.png',
                    'Aggregated Message Reception Delay with `pub0sub`'))


if __name__ == '__main__':
    main()
