#!/usr/bin/env python3
# -*- coding=utf-8 -*-

import re
import time
from concurrent.futures.thread import ThreadPoolExecutor
from pathlib import Path


def replace_not_null(content: str):
    p = re.compile(r"\s?(?P<name>[\w_]+(\(\))?)\s?!=\s?null\s?", re.VERBOSE)
    res = p.subn(r'Objects.nonNull(\g<name>)', content)
    return res

d = replace_not_null("if (ooo() != null)")
print(d)
print(replace_not_null(d[0]))

print(re.subn('(xx)+', '', 'abcdab'))
print(re.subn('(ab)+', '', 'abcdab'))