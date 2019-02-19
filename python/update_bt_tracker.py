#!/usr/bin/env python3
# -*- coding=utf-8 -*-
import re
import subprocess
import os
from pathlib import Path
from pathlib import PosixPath

pkgs = ['requests']
for package in pkgs:
    try:
        import package
    except ImportError as e:
        subprocess.check_call(["python3", '-m', 'pip', 'install', package])

import requests


def get_update_data(path: str, bt_tracker: str):
    with open(path, "r", encoding="utf-8") as f:
        data = f.read()
        updated = re.sub(r"bt-tracker=([\w,:/.\d-]+)", bt_tracker, str(data))
        print(updated)
        return updated


trackerlist_url = "https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all.txt"

r = requests.get(trackerlist_url, allow_redirects=True)
if r.content is None:
    print("respone is empty!")
    exit()
content = r.content.replace(str.encode("\n\n"), str.encode(","), -1)
bt_tracker = "bt-tracker={}".format(content.decode("utf-8"))
current_dir = PosixPath(".")
aria2_conf_path = Path(current_dir.home())
aria2_conf_path = aria2_conf_path.joinpath(".aria2/aria2.conf")
if not aria2_conf_path.exists():
    print("aria2.conf did not exists")
    exit()

updated_data = get_update_data(aria2_conf_path, bt_tracker)
try:
    back_file_path = Path(str(aria2_conf_path.absolute()
                              ).replace(".conf", ".conf.bak"))
    if back_file_path.exists:
        back_file_path.unlink()
        aria2_conf_path.rename(back_file_path)
except Exception as e:
    print(e)
    aria2_conf_path.rename(back_file_path)

aria2_conf_path.open("w", encoding="utf-8").write(updated_data)
