#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import time
from concurrent.futures.thread import ThreadPoolExecutor
from pathlib import Path


def find(file: Path):
    # TODO: need implementation
    pass


def main(dir_path: Path):
    start_time = time.time()
    executor = ThreadPoolExecutor(max_workers=10)
    for path in dir_path.rglob("*.java"):
        executor.submit(find, path)
    end_time = time.time()
    print("estimate time: %d".format(end_time - start_time))


if __name__ == "__main__":
    main(Path.cwd())
