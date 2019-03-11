#!/bin/bash
# add script to crontab
# reboot server when swap less then 20mb.


COMMAND='59 23 * * * root ~/Rorscipt/shell/monitor_swap.sh'

EXISTS=`cat /etc/crontab | grep monitor_swap.sh`

if [ $EXISTS -z ];
then
	echo -e $COMMAND >> /etc/crontab
fi


USAGE=`free -m | awk 'NR==3{print $4}'`
if [ $USAGE -lt 20 ]; then
	reboot
fi
