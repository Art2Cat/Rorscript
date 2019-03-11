#!/bin/bash
# add script to crontab
# reboot server when swap less then 20mb.

USAGE=`free -m | awk 'NR==3{print $4}'`
if [ $USAGE -lt 20 ]; then
	reboot
fi
