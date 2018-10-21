#!/bin/bash

function start_aria2 {
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

function update_bt_tracker {
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
    # print(r.content)
    bt_tracker = "bt-tracker={}".format(content.decode("utf-8"))
    print(bt_tracker)
    aria2_conf_path = os.path.join("~/.aria2/", "aria2.conf")
    if os.path.exists(aria2_conf_path):
        updated_data = get_update_data(aria2_conf_path, bt_tracker)
        with open(aria2_conf_path, "w", encoding="utf-8") as writer:
            writer.write(updated_data)

    EOF

    chmod 755 update_bt_tracker.py

    python3 ./update_bt_tracker.py

    rm -f ./update_bt_tracker.py
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
		    # update bt tracker before start aria2
		    update_bt_tracker ;;
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
