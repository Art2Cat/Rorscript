#!/usr/bin/env python3
# -*- coding=utf-8 -*-
import os
import platform
import re
import time


def find_from_table_names(content: str):
    result = set([])
    exclude = ["as", "is"]
    # print(all_content)
    for m in re.finditer(r'((from|FROM)\s([a-zA-z_]+)\s)', content):
        if m.group(3) not in exclude:
            result.add(m.group(3))
    return result


def find_insert_table_names(content: str):
    result = set([])
    for m in re.finditer(r'((insert)\s(into)\s([a-zA-Z_]+)\s?)',
                         content, re.IGNORECASE):
        result.add(m.group(4))
        print(m.groups())
    return result


def find_join_table_names(content: str):
    result = set([])
    for m in re.finditer(r'((join|JOIN)\s([a-zA-Z_]+)\s)', content):
        result.add(m.group(3))
    return result


def find_update_table_names(content: str):
    result = set([])
    for m in re.finditer(r'((update|UPDATE)\s([a-zA-Z_]+)\s)', content):
        result.add(m.group(3))
    return result


def find_table_names(file_path: str):
    result = dict({})
    with open(file_path, 'r') as f:

        all_content = f.read()
        result["from"] = find_from_table_names(content=all_content)
        result["insert"] = find_insert_table_names(content=all_content)
        result["join"] = find_join_table_names(content=all_content)
        result["update"] = find_update_table_names(content=all_content)
    return result


def replace_tables(data: str, tables: set, ptn: str, value: str):
    for v in tables:
        ptn = ptn.format(v)
        print(ptn)
        value = value.format(v)
        print(value)
        data = re.sub(ptn, value, data)
    # print(data)
    return data


def replace(input_file, output, tables):
    with open(input_file, 'r') as reader:

        replaced = ""
        data = reader.read()
        with open(output, 'w') as writer:
            from_tables = tables.get("from")
            for i, v in enumerate(from_tables):
                ptn = r"(from|FROM)\s{}".format(v)
                val = "FROM dbo.{}".format(v)
                if i == 0:
                    replaced = re.sub(ptn, val, data)
                else:
                    replaced = re.sub(ptn, val, replaced)
            insert_tables = tables.get("insert")
            print(len(insert_tables))
            # ptn = r"((insert|INSERT|Insert)\s(into|INTO)\s{}\s?)"
            # val = "INSERT INTO dbo.{} "
            # replaced = replace_tables(replaced, insert_tables, ptn, val)
            for v in insert_tables:
                ptn = r"(([insertINSERT]+)\s(into|INTO)\s{}\s?)".format(v)
                val = "INSERT INTO dbo.{} ".format(v)
                replaced = re.sub(ptn, val, replaced, re.IGNORECASE)

            join_tables = tables.get("join")
            for v in join_tables:
                ptn = r"(join|JOIN)\s{}".format(v)
                val = "JOIN dbo." + v
                replaced = re.sub(ptn, val, replaced)

            update_tables = tables.get("update")
            # ptn = "(update|UPDATE) {}"
            # val = "update dbo.{}"
            # replaced = replace_tables(replaced, update_tables, ptn, val)
            for v in update_tables:
                ptn = r"(update|UPDATE)\s{}".format(v)
                val = "update dbo.{}".format(v)
                replaced = re.sub(ptn, val, replaced)

            writer.write(replaced)


def main(path: str):
    tables: dict
    input_file: str
    output_file: str
    for root, dirs, files in os.walk(path):
        for file in files:
            file_name = os.path.basename(file)
            file_path = os.path.join(root, file_name)
            if os.path.isfile(file_path) and ".xml" in file_path:
                input_file = file_path
                output_file = os.path.join(
                    root, "{} - output.xml".format(file_name))
                print(file_path)
                tables = find_table_names(file_path)
                replace(input_file, output_file, tables)
                tables.clear()


if __name__ == "__main__":
    start_time = time.time()
    path = os.path.abspath('.')
    print(path)
    main(path)
    end_time = time.time()
    print(end_time - start_time)
