#!/bin/bash

function start_aria2 {
	nohup aria2c -c &>~/.aria2/aria2.log &
}

function stop_aria2 {
	kill -9 $(ps -ef|grep "aria2c"|grep -v "grep"|awk '{print $2}')
}

function showlog {
	tail -f ~/.aria2/aria2.log
}

function clearlog {
	echo "" > ~/.aria2/aria2.log
}

function menu {
	clear
	echo
	echo -e "\t\t\tAria2 Admin Menu\n"
	echo -e "\t1. Start aria2"
	echo -e "\t2. Stop aria2"
	echo -e "\t3. Monitor aria2 log"
	echo -e "\t3. clear aria2 log"
	echo -e "\t0. Exit program\n\n"
	echo -en "\t\tEnter option: "
	read -n 1 option
}

while [ 1 ]
do
	menu
	case $ option in
		0)
			break ;;
		1)
			start_aria2 ;;
		2)
			stop_aria2 ;;
		3)
			showlog ;;
		4)
			clearlog ;;
		*)
			clear
			echo "Sorry, wrong selection";;
	esac
	echo -en "\n\n\t\t\tHit any key to continue"
	read -n 1 line
done
clear
