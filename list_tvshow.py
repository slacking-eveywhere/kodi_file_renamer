#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from pathlib import Path


class TVShowList:
    def __init__(self, tvshow_path):
        self.tvshow_path = tvshow_path

    def __iter__(self):
        for tvshow in sorted(os.listdir(self.tvshow_path)):
            tvshow_path = Path(os.path.join(self.tvshow_path, tvshow))
            if not tvshow_path.name.startswith("."):
                yield tvshow_path.name
