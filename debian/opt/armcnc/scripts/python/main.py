#!/usr/bin/env python
# -*- coding: utf-8 -*-

def __init__(self, **words):
    print(len(words), " words passed")
    for w in words:
        print("%s: %s" % (w, words[w]))
