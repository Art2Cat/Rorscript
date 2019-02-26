#!/usr/bin/env python3
# -*- coding=utf-8 -*-
import re
import time
from concurrent.futures.thread import ThreadPoolExecutor
from pathlib import Path


def add_import(content: str):
    p = re.compile(r"import\s(?P<name>(static)?\s?[\w.;]+)", re.VERBOSE)
    return p.subn(r"import java.util.Objects;\nimport \g<name>", content, 1)


def replace_is_null(content: str):
    p = re.compile(
        r"\s?(?P<name>[\w.\d\[\]_]+(\(\))?)\s?==\s?null\s?", re.VERBOSE)
    res = p.subn(r'Objects.isNull(\g<name>)', content)
    return res


def replace_not_null(content: str):
    p = re.compile(
        r"\s?(?P<name>[\w.\d\[\]_]+(\(\))?)\s?!=\s?null\s?", re.VERBOSE)
    res = p.subn(r'Objects.nonNull(\g<name>)', content)
    return res


def replace(file_path: Path):
    state = set()
    rss = file_path.read_text(encoding="utf-8")
    rss = replace_is_null(rss)
    state.add(rss[1])
    rss = replace_not_null(rss[0])
    state.add(rss[1])
    if 0 not in state or len(state) == 2:
        if "import java.util.Objects" not in rss[0]:
            rss = add_import(rss[0])
        new = file_path.parent.joinpath(file_path.name + ".out")
        with open(str(new), "w", encoding="utf-8") as f:
            f.write(rss[0])
        if new.exists():
            file_path.rename(file_path.parent.joinpath(
                file_path.name + ".old"))
            new.rename(file_path)


def main(dir_path: Path):
    start_time = time.time()
    executor = ThreadPoolExecutor(max_workers=10)
    for path in dir_path.rglob("*.java"):
        executor.submit(replace, path)
    end_time = time.time()
    print(end_time - start_time)


if __name__ == "__main__":
    main(Path.cwd())
