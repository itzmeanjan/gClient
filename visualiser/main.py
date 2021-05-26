#!/usr/bin/python3

from argparse import ArgumentParser
from os.path import abspath, isabs, isdir, join
from os import walk
from typing import Dict, List
from re import compile as re_compile


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
    accumulate_data(found[0], bucket)
    print(bucket)


if __name__ == '__main__':
    main()
