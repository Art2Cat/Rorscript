#!/usr/bin/env python3
# -*- coding=utf-8 -*-

import re
import time
from concurrent.futures.thread import ThreadPoolExecutor
from pathlib import Path


def add_import(content: str):
    p = re.compile(r"import\s(?P<name>(java.util.)?[\w.;]+)", re.VERBOSE)
    return p.sub(r"import \g<name>\nimport java.util.Objects;", content, 1)


def replace_is_null(content: str):
    p = re.compile(r"\s?(?P<name>[\w_]+)\s?==\s?null\s?", re.VERBOSE)
    res = p.sub(r'Objects.isNull(\g<name>)', content)
    return res


def replace_not_null(content: str):
    p = re.compile(r"\s?(?P<name>[\w_]+)\s?!=\s?null\s?", re.VERBOSE)
    res = (p.sub(r'Objects.nonNull(\g<name>)', content))
    return res


def replace(file_path: Path):
    rss = file_path.read_text(encoding="utf-8")
    rss = replace_is_null(rss)
    rss = replace_not_null(rss)
    rss = add_import(rss)
    with open(str(file_path.parent.joinpath(file_path.name + ".out")), "w", encoding="utf-8") as f:
        f.write(rss)


def main(dir_path: Path):
    start_time = time.time()
    executor = ThreadPoolExecutor(max_workers=10)
    for path in dir_path.rglob("*.java"):
        executor.submit(replace, path)
    end_time = time.time()
    print(end_time - start_time)


if __name__ == "__main__":
    main(Path.cwd())
