#!/bin/bash


function start_aria2 {
	# update bt tracker before start aria2
	python3 ~/Rorscript/python/update_bt_tracker.py &

	nohup aria2c --enable-rpc --rpc-listen-all -c &>~/.aria2/aria2.log &
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
