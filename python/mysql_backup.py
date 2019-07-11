# -*- coding: utf-8 -*-

from pathlib import Path

import MySQLdb
import os
import datetime
# require mysqlclient

db = 'test'

con = MySQLdb.connect(host='localhost', user='root', passwd='password', db=db, charset='utf8mb4')
cur = con.cursor()

cur.execute("SHOW TABLES")
data = ""
tables = []
for table in cur.fetchall():
    tables.append(table[0])

for table in tables:
    data += "DROP TABLE IF EXISTS `" + str(table) + "`;"

    cur.execute("SHOW CREATE TABLE `" + str(table) + "`;")
    data += "\n" + str(cur.fetchone()[1]) + ";\n\n"

    cur.execute("SELECT * FROM `" + str(table) + "`;")
    for row in cur.fetchall():
        data += "INSERT INTO `" + str(table) + "` VALUES("
        first = True
        for field in row:
            if not first:
                data += ', '
            data += '"' + str(field) + '"'
            first = False

        data += ");\n"
    data += "\n\n"

now = datetime.datetime.now()

filename = "backup_" + db + "_" + now.strftime("%Y-%m-%d_%H_%M_%S")
path = Path(os.getcwd(), filename)
path = path.with_suffix(".sql")

with path.open("w", encoding="utf-8") as f:
    f.write(data)
