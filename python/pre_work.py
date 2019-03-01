#!/usr/bin/env python3
# -*- coding=utf-8 -*-
import re
from pathlib import Path


def replace_java_name(content: str):
    p = re.compile(r"(?P<java>    ([a-z.]+)  )", re.VERBOSE)
    res = p.subn(r'package: \g<java>\n', content)
    return res


def replace(file_path: Path):
    rss = file_path.read_text(encoding="utf-8")
    rss = replace_java_name(rss)
    new = file_path.parent.joinpath(file_path.name + ".out")
    with open(str(new), "w", encoding="utf-8") as f:
        f.write(rss[0])
    if new.exists():
        file_path.rename(file_path.parent.joinpath(
            file_path.name + ".old"))
        new.rename(file_path)


if __name__ == "__main__":
    replace(Path.cwd().joinpath("deprecated.txt"))
