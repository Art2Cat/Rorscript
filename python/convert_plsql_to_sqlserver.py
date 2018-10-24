#!/usr/bin/env python3
# -*- coding=utf-8 -*-
import os
import platform
import re
import time


def replace(input_file, output, tables):
    with open(input_file, 'r') as reader:
        replaced = ""
        data = reader.read()
        with open(output, 'w') as writer:
            replaced = data.encode(encoding="utf-8")
            replaced = re.sub(
                r"(from|FROM)\s(?P<table>[a-zA-Z_]+)", r"FROM dbo.\g<table>", replaced)
            replaced = re.sub(
                r"(([insertINSERT]+)\s(into|INTO)\s(?P<table>[a-zA-Z_]+)\s?)", r"INSERT INTO dbo.\g<table> ", replaced)
            replaced = re.sub(
                r"(join|JOIN)\s(?P<table>[a-zA-Z_]+)", r"JOIN dbo.\g<table>", replaced)
            replaced = re.sub(r"(update|UPDATE)\s(?P<table>[a-zA-Z_]+)",
                              r"update dbo.\g<table>", replaced)
            writer.write(replaced)


def main(path: str):
    input_file: str
    output_file: str
    for root, _, files in os.walk(path):
        for file in files:
            file_name = os.path.basename(file)
            file_path = os.path.join(root, file_name)
            if os.path.isfile(file_path) and ".xml" in file_path:
                input_file = file_path
                output_file = os.path.join(
                    root, "{} - output.xml".format(file_name))
                print(file_path)
                replace(input_file, output_file, tables)


if __name__ == "__main__":
    start_time = time.time()
    path = os.path.abspath('.')
    print(path)
    main(path)
    end_time = time.time()
    print(end_time - start_time)
