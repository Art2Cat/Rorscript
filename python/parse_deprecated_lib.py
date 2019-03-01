#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re
from pathlib import Path


class Lib(object):
    _name: str

    _classes: dict

    def __init__(self):
        self._name = None
        self._classes = dict()

    @property
    def name(self):
        return self._name

    @name.setter
    def name(self, name: str):
        self._name = name

    @property
    def classes(self):
        return self._classes

    def __iter__(self):
        return (i for i in (self.name, self.classes))

    def __repr__(self):
        class_name = type(self).__name__
        return '{}({!r}, {!r})'.format(class_name, *self)


class DeprecatedClass():
    _name: str
    _full_name: str
    _methods: set

    def __init__(self):
        self._name = None
        self._full_name = None
        self._methods = set()

    @property
    def name(self):
        return self._name

    @name.setter
    def name(self, name: str):
        self._name = name

    @property
    def full_name(self):
        return self._full_name

    @full_name.setter
    def full_name(self, name: str):
        self._full_name = name

    @property
    def methods(self):
        return self._methods

    @methods.setter
    def methods(self, methods: set):
        self._methods = methods

    def __iter__(self):
        return (i for i in (self.name, self.full_name, self.methods))

    def __repr__(self):
        class_name = type(self).__name__
        return '{}({!r}, {!r})'.format(class_name, *self)


def get_libname(content: str):
    match = re.search(r":\s([\w.:_-]+)", content)
    if match is not None:
        return match.group(1)


def get_package_name(content: str):
    match = re.search(r"([^@][a-z._-]+)", content)
    if match is not None:
        return match.group(0)


def get_java_file_name(content: str):
    match = re.search(r"([\w]+\.java)", content)
    if match is not None:
        return match.group(0)


def get_method_name(content: str):
    match = re.search(r"([^@][a-zA-Z]+)\(", content)
    if match is not None:
        return match.group(1)


def parse(dir_path: Path):
    lib_file = dir_path.joinpath("deprecated.txt")
    content = lib_file.read_text("utf-8")
    libs = set()
    for mvn in content.split("Maven"):
        lib = Lib()
        lines = mvn.split("\n")
        lib.name = get_libname(lines[0])
        for line in lines:
            clz = get_java_file_name(line)
            if clz is not None:
                dcls = DeprecatedClass()
                dcls.name = clz
                lib.classes[clz] = dcls

        libs.add(lib)

    for l in libs:
        print(l)


if __name__ == "__main__":
    parse(Path.cwd())
