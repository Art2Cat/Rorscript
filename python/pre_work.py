#!/usr/bin/env python3
# -*- coding=utf-8 -*-
import re
from pathlib import Path
from collections import UserString
import xml.etree.ElementTree as ET
from xml.dom import minidom


def remove_usage_found(content: str):
    p = re.compile(r"(\([\d]+\susage\sfound\))", re.VERBOSE)
    res = p.sub(r"", content)
    return res


def remove_usages_found(content: str):
    p = re.compile(r"(\([\d]+\susages\sfound\))", re.VERBOSE)
    res = p.sub(r"", content)
    return res


def remove_deprecated_annotation(content: str):
    p = re.compile(r"(@Deprecated)", re.VERBOSE)
    res = p.sub(r"", content)
    return res


def remove_comment(content: str):
    p = re.compile(r"(//\s?[\w.:?(\-)\s]+)", re.VERBOSE)
    res = p.sub(r"", content)
    return res


def remove_empty_line(content: str):
    new_rss = UserString("")
    for line in content.split("\n"):
        if not line.strip():
            continue
        if "mvn" in line:
            new_rss += line.strip()
        else:
            new_rss += line.rstrip()
        new_rss += "\n"
    return str(new_rss)


def concatenate(org, new: str):
    org += new


def add_lib_tag(content: str):
    p = re.compile(r"Maven:\s(?P<lib>[\w.:_-]+)", re.VERBOSE)
    res = p.sub(r"mvn: \g<lib>\n", content)
    return res


def add_package_tag(content: str):
    new_ss = UserString("")
    for line in content.split("\n"):
        state = set()
        if re.search(r"(mvn:)", line) is not None:
            state.add(1)
        if re.search(r"(class:)", line) is not None:
            state.add(1)
        if re.search(r"(mtd:)", line) is not None:
            state.add(1)

        if len(state) == 0:
            p = re.compile(r"\s{12}(?P<package>[a-z0-9.:]*)", re.VERBOSE)
            res = p.sub(r"   package: \g<package>\n", line)
            new_ss += res
        else:
            new_ss += line
            new_ss += "\n"
    return str(new_ss)


def add_class_tag(content: str):
    p = re.compile(r"(?P<java>[\w]+\.java)", re.VERBOSE)
    res = p.sub(r"\n        class: \g<java>\n", content)
    return res


def add_method_tag(content: str):
    p = re.compile(r"(?P<mtd>[\w]+\([\w<>?,\s.\[\]]*\))", re.VERBOSE)
    res = p.sub(r"\n            mtd: \g<mtd>\n", content)
    return res


def wash_data(file_path: Path):
    rss = file_path.read_text(encoding="utf-8")
    rss = remove_usage_found(rss.strip())
    rss = remove_usages_found(rss)
    rss = remove_deprecated_annotation(rss)
    rss = remove_comment(rss)
    rss = add_lib_tag(rss)
    rss = add_class_tag(rss)
    rss = add_method_tag(rss)
    rss = remove_empty_line(rss)
    rss = add_package_tag(rss)
    return rss


def get_mvn_name(content: str):
    match = re.search(r"mvn:\s([\w.:_-]+)", content)
    if match is not None:
        return match.group(1)


def get_package_name(content: str):
    match = re.search(r"package:\s([a-z0-9.]+)", content)
    if match is not None:
        return match.group(1)


def get_class_name(content: str):
    match = re.search(r"class:\s([\w]+)\.java", content)
    if match is not None:
        return match.group(1)


def get_method_name(content: str):
    match = re.search(r"mtd:\s([a-zA-Z]+)\(", content)
    if match is not None:
        return match.group(1)


def prettify(elem):
    """
    Return a pretty-printed XML string for the Element.
    """
    rough_string = ET.tostring(elem, "utf-8")
    reparsed = minidom.parseString(rough_string)
    return reparsed.toprettyxml(indent="    ")


def to_xml(content: str):
    data = ET.Element("data")
    mvn = None
    pkg = None
    clz = None
    for line in content.split("\n"):

        if re.search(r"(mvn:)", line) is not None:
            mvn = ET.SubElement(data, 'mvn')
            mvn.set("name", get_mvn_name(line))
        if re.search(r"(package:)", line) is not None:
            if mvn is not None:
                pkg = ET.SubElement(mvn, "pkg")
                pkg.set("name", get_package_name(line))
        if re.search(r"(class:)", line) is not None:
            if pkg is not None:
                clz = ET.SubElement(pkg, 'class')
                clz.set("name", get_class_name(line))
        if re.search(r"(mtd:)", line) is not None:
            if clz is not None:
                name = get_method_name(line)
                if name is not None and name[0].islower():
                    mtd = ET.SubElement(clz, "mtd")
                    mtd.text = name

    # mydata = ET.tostring(data)
    with open(str(Path.cwd().joinpath("deprecated.xml")), "w") as f:
        f.write(prettify(data))


if __name__ == "__main__":
    washed = wash_data(Path.cwd().joinpath("deprecated.txt"))
    to_xml(washed)