#!/bin/bash
# reboot server when swap less then 20mb.
while true;
do

	USAGE=`free -m | awk 'NR==3{print $4}'`
	if [ $USAGE -lt 20 ]; then
		reboot
	fi
done
