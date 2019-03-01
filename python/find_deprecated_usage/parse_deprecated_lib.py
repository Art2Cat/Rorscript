#!/usr/bin/env python3
# -*- coding: utf-8 -*-
from pathlib import Path
from xml.dom import minidom

from find_deprecated_usage.pre_work import wash_data, to_xml


class Lib(object):
    _name: str
    _packages: set

    def __init__(self):
        self._name = None
        self._packages = set()

    @property
    def name(self):
        return self._name

    @name.setter
    def name(self, name: str):
        self._name = name

    @property
    def packages(self):
        return self._packages

    def __iter__(self):
        return (i for i in (self.name, self.packages))

    def __repr__(self):
        class_name = type(self).__name__
        return '{}({!r}, {!r})'.format(class_name, *self)


class Package:
    _name: str
    _classes: set

    def __init__(self):
        self._name = None
        self._classes = set()

    @property
    def name(self):
        return self._name

    @name.setter
    def name(self, name: str):
        self._name = name

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

    def __iter__(self):
        return (i for i in (self.name, self.full_name, self.methods))

    def __repr__(self):
        class_name = type(self).__name__
        return '{}({!r}, {!r})'.format(class_name, *self)


def parse():
    washed = wash_data(Path.cwd().joinpath("export.txt"))
    to_xml(washed)

    mydoc = minidom.parse('deprecated.xml')

    items = mydoc.getElementsByTagName("mvn")
    libs = set()
    for i in items:
        lib = Lib()
        lib.name = (i.attributes['name'].value)
        pkgs = i.getElementsByTagName("pkg")
        for p in pkgs:
            pkg = Package()
            pkg.name = p.attributes['name'].value
            clzs = p.getElementsByTagName("class")
            for c in clzs:
                clz = DeprecatedClass()
                clz.name = c.attributes["name"].value
                clz.full_name = pkg.name + "." + clz.name
                mtds = c.getElementsByTagName("mtd")
                for m in mtds:
                    if m.firstChild is not None:
                        clz.methods.add(m.firstChild.nodeValue)
                pkg.classes.add(clz)
            lib.packages.add(pkg)
        libs.add(lib)

    return libs


if __name__ == "__main__":
    parse()
