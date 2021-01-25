#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from pathlib import Path


class TVShowList:
    def __init__(self, tvshows_path):
        self.tvshows_path = tvshows_path

    def __iter__(self):
        for tvshow in sorted(os.listdir(self.tvshows_path)):
            tvshow_path = Path(os.path.join(self.tvshows_path, tvshow))
            if tvshow_path.name.startswith("."):
                continue

            yield tvshow_path.name


class TVShowEpisodesList:
    def __init__(self, tvshow_path):
        self.tvshow_path = tvshow_path

    def __iter__(self):
        for episode in sorted(os.listdir(self.tvshow_path)):
            episode_path = Path(os.path.join(self.tvshow_path, episode))
            if episode_path.name.startswith("."):
                continue

            yield episode_path
