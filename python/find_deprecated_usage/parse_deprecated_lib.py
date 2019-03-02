#!/usr/bin/env python3
# -*- coding: utf-8 -*-
from pathlib import Path
from xml.dom import minidom

from find_deprecated_usage.pre_work import wash_data, to_xml


class DeprecatedClass:
    """
    :type _name: str
    :type _full_name: str
    :type _pkg_name: str
    :type _lib_name: str
    :type _methods: set
    :type _usage: int

    """
    __slots__ = (
        "_name", "_full_name", "_pkg_name", "_lib_name", "_methods", "_usage"
    )

    def __init__(self):
        self._name = None
        self._full_name = None
        self._pkg_name = None
        self._lib_name = None
        self._methods = set()
        self._usage = 0

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
    def pkg_name(self):
        return self._pkg_name

    @pkg_name.setter
    def pkg_name(self, name: str):
        self._pkg_name = name

    @property
    def lib_name(self):
        return self._lib_name

    @lib_name.setter
    def lib_name(self, name: str):
        self._lib_name = name

    @property
    def methods(self):
        return self._methods

    @property
    def usage(self):
        return self._lib_name

    @usage.setter
    def usage(self, u: int):
        self._usage = u

    def __iter__(self):
        return (i for i in (self.name, self.full_name, self.pkg_name, self.lib_name, self.methods, self.usage))

    def __repr__(self):
        class_name = type(self).__name__
        return "{}({!r}, {!r}, {!r}), {!r}, {!r}, {!r}".format(class_name, *self)


def parse():
    washed = wash_data(Path.cwd().joinpath("export.txt"))
    to_xml(washed)

    mydoc = minidom.parse("deprecated.xml")

    items = mydoc.getElementsByTagName("mvn")
    libs = set()
    for i in items:
        lib_name = i.attributes["name"].value
        pkgs = i.getElementsByTagName("pkg")
        for p in pkgs:
            pkg_name = p.attributes["name"].value
            clzs = p.getElementsByTagName("class")
            for c in clzs:
                clz = DeprecatedClass()
                clz.name = c.attributes["name"].value
                clz.full_name = "{}.{}".format(pkg_name, clz.name)
                clz.pkg_name = pkg_name
                clz.lib_name = lib_name
                mtds = c.getElementsByTagName("mtd")
                for m in mtds:
                    if m.firstChild is not None:
                        clz.methods.add(m.firstChild.nodeValue)
                libs.add(clz)

    return libs


if __name__ == "__main__":
    parse()
