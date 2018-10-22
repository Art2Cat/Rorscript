#!/bin/bash

cat << EOF > update_bt_tracker.py

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
bt_tracker = "bt-tracker={}".format(content.decode("utf-8"))
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

EOF

chmod 755 update_bt_tracker.py

function start_aria2 {
	# update bt tracker before start aria2
	python3 ./update_bt_tracker.py &

	rm -f ./update_bt_tracker.py
	nohup aria2c -c &>~/.aria2/aria2.log &
}

function stop_aria2 {
	kill -9 $(ps -ef|grep "aria2c"|grep -v "grep"|awk '{print $2}')
}

function show_log {
	tail -f ~/.aria2/aria2.log
}

function clear_log {
	echo "" > ~/.aria2/aria2.log
}


function menu {
	clear
	echo
	echo -e "\t\t\tAria2 Admin Menu\n"
	echo -e "\t1. Start aria2"
	echo -e "\t2. Stop aria2"
	echo -e "\t3. Monitor aria2 log"
	echo -e "\t4. clear aria2 log"
	echo -e "\t0. Exit program\n\n"
	echo -en "\t\tEnter option: "
	read -n 1 option
}

while [ 1 ]
do
	menu
	case $option in
		0)
			break ;;
		1)
			start_aria2 ;;
		2)
			stop_aria2 ;;
		3)
			show_log ;;
		4)
			clear_log ;;
		*)
			clear
			echo "Sorry, wrong selection";;
	esac
	echo -en "\n\n\t\t\tHit any key to continue"
	read -n 1 line
done
clear
