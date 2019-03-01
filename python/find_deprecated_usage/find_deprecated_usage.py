#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re
from concurrent.futures.thread import ThreadPoolExecutor
from datetime import datetime
from pathlib import Path

from find_deprecated_usage.parse_deprecated_lib import parse


def find(file: Path, libs: set):
    # TODO: write data to excel file
    need_check_mtds = set()
    potential_libs = dict()
    rss = file.read_text(encoding="utf-8")
    for s in rss.split("\n"):
        for l in libs:
            for p in l.packages:
                for c in p.classes:
                    r = "({})".format(c.full_name)
                    if re.search(r, s) is not None:
                        potential_libs[c.name] = l.name
                        need_check_mtds.update(c.methods)

    usage = 0
    for s in rss.split("\n"):
        for m in need_check_mtds:
            r = "({})".format(m)
            if re.search(r, s) is not None:
                usage += 1

    print("{} deprecated method usage {}".format(file.name, usage))
    for k, v in potential_libs.items():
        print("used class: {} in lib: {}".format(k, v))


def main(dir_path: Path):
    start_time = datetime.now()
    libs = parse()
    executor = ThreadPoolExecutor(max_workers=10)
    for path in dir_path.rglob("*.java"):
        executor.submit(find, path, libs)
        # find(path, libs)
    end_time = datetime.now()
    print("estimate time: {}".format(end_time - start_time))


if __name__ == "__main__":
    main(Path.cwd())
