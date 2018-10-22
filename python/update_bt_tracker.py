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


def get_update_data(path: str, bt_tracker: str):
    with open(path, "r", encoding="utf-8") as f:
        data = f.read()
        updated = re.sub(r"bt-tracker=([\w,:/.\d-]+)", bt_tracker, str(data))
        print(updated)
        return updated

import requests

trackerlist_url = "https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all.txt"

r = requests.get(trackerlist_url, allow_redirects=True)
content = r.content.replace(str.encode("\n\n"), str.encode(","), -1)
# print(r.content)
bt_tracker = "bt-tracker={}".format(content.decode("utf-8"))
print(bt_tracker)
aria2_conf_path = os.path.join("~/.aria2/", "aria2.conf")
if not os.path.exists(aria2_conf_path):
    print("aria2.conf did not exists: " + aria2_conf_path)
    exit()

updated_data = get_update_data(aria2_conf_path, bt_tracker)
try:
    back_file_path = aria2_conf_path + ".bak"
    if os.path.exists(back_file_path):
        os.remove(back_file_path)
        os.rename(aria2_conf_path, back_file_path)
except Exception as e:
    print(e)
    exit()
with open(aria2_conf_path, "w", encoding="utf-8") as writer:
    writer.write(updated_data)
