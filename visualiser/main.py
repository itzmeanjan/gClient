#!/usr/bin/python3

from argparse import ArgumentParser
from os.path import abspath, isabs, isdir, join
from os import walk
from typing import Dict, List, Pattern, Tuple
from re import compile as re_compile
from math import ceil
from matplotlib import pyplot as plt
import seaborn as sns
from humanize import precisedelta
from datetime import timedelta


def visualise(data: Dict[Tuple[int, int], int], sink: str, xlabel: str, ylabel: str, title: str, subtitle: str) -> bool:
    if not data:
        return False

    try:
        tmp = [(f'{start} - {end}ms', count)
               for (start, end), count in data.items()]
        x = [t[0] for t in tmp]
        y = [t[1] for t in tmp]

        with plt.style.context('seaborn-darkgrid'):
            fig = plt.Figure(
                figsize=(18, 9),
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

            fig.gca().set_xlabel(xlabel, labelpad=12)
            fig.gca().set_ylabel(ylabel, labelpad=12)
            fig.gca().set_title(subtitle, pad=6, fontsize=15)
            fig.suptitle(title, fontsize=20, y=1)

            fig.savefig(
                sink,
                bbox_inches='tight',
                pad_inches=.5)
            plt.close(fig)

        return True
    except Exception as e:
        print(f'Error : {e}')
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
    if max_delay - min_delay < 2*slot:
        return [(min_delay, max_delay)]

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


def accumulate_data(file: str, bucket: Dict[int, int], record_duration: Dict[str, int]):
    with open(file) as fd:
        while True:
            ln = fd.readline()
            if not ln:
                break

            ts = [i.strip() for i in ln.split(';')][:2]
            sent = int(ts[0])
            received = int(ts[1])

            if sent < record_duration['start']:
                record_duration['start'] = sent
            if received > record_duration['end']:
                record_duration['end'] = received

            diff = received - sent
            if diff not in bucket:
                bucket[diff] = 1
            else:
                bucket[diff] = bucket[diff] + 1

    return


def find_files(dir: str, reg: Pattern) -> List[str]:
    found = []
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

    pub_reg = re_compile(r'^(log_pub_\d+\.csv)$')
    sub_reg = re_compile(r'^(log_sub_\d+\.csv)$')

    pub_log = find_files(target_dir, pub_reg)
    if not pub_log:
        print('Directory walk found no publisher log !')
        return

    sub_log = find_files(target_dir, sub_reg)
    if not sub_log:
        print('Directory walk found no subscriber log !')
        return

    pub_record_duration = {'start': 1 << 64 - 1, 'end': 0}
    pub_bucket = {}
    for file in pub_log:
        accumulate_data(file, pub_bucket, pub_record_duration)

    sub_record_duration = {'start': 1 << 64 - 1, 'end': 0}
    sub_bucket = {}
    for file in sub_log:
        accumulate_data(file, sub_bucket, sub_record_duration)

    pub_dt = timedelta(
        milliseconds=pub_record_duration['end'] - pub_record_duration['start'])
    sub_dt = timedelta(
        milliseconds=sub_record_duration['end'] - sub_record_duration['start'])

    pub_slots = splitted_delay_spectrum(pub_bucket.keys(), args.slot)
    status = visualise(aggregated_count_by_slot(pub_slots, pub_bucket),
                       'pub_out.png',
                       'Message Received After Delay',
                       '#-of Messages',
                       'Aggregated Message Reception Delay with `pub0sub`',
                       f'Recorded for {precisedelta(pub_dt)}')
    if status:
        print(f'Publisher visualisation ✅')

    sub_slots = splitted_delay_spectrum(sub_bucket.keys(), args.slot)
    status = visualise(aggregated_count_by_slot(sub_slots, sub_bucket),
                       'sub_out.png',
                       'Message Sending Latency',
                       '#-of Messages',
                       'Aggregated Message Sending Latency with `pub0sub`',
                       f'Recorded for {precisedelta(sub_dt)}')
    if status:
        print(f'Subscriber visualisation ✅')


if __name__ == '__main__':
    main()
