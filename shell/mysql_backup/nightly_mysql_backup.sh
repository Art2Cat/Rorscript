#!/bin/bash
# crontab -e 0 4 * * * ~/nightly_mysql_backup.sh
mysqldump --defaults-extra-file=./client.cnf  db_scheme > ~/db_backup/db_$( date +"%Y_%m_%d" ).sql

find ~/db_backup/ -type f -mtime +7 -name '*.sql' -execdir rm -- '{}' +

