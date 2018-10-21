#!/usr/bin/env python3
# -*- coding=utf-8 -*-
import re
import subprocess
import os

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
content = r.content.replace(str.encode("\n\n"), str.encode(","), -1)
# print(r.content)
bt_tracker = "bt-tracker={}".format(content.decode("utf-8"))
print(bt_tracker)
aria2_conf_path = os.path.join("~/.aria2/", "aria2.conf")
if os.path.exists(aria2_conf_path):
    updated_data = get_update_data(aria2_conf_path, bt_tracker)
    with open(aria2_conf_path, "w", encoding="utf-8") as writer:
        writer.write(updated_data)
