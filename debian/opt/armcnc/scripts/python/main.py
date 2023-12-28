#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# import脚本，例如换刀、对刀、夹具的Python脚本

def __init__(self, **words):
    print(len(words), " words passed")
    for w in words:
        print("%s: %s" % (w, words[w]))
