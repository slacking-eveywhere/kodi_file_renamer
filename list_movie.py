#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os

MOVIE_FILE_TYPES = [".mkv", ".avi", ".mp4", ".iso", ".m2ts"]


class MovieList:
    def __init__(self, movie_path):
        self.movie_path = movie_path

    def __iter__(self):
        for movie in sorted(os.listdir(self.movie_path)):
            movie_filename = os.path.join(self.movie_path, movie)
            movie_name, extension = os.path.splitext(movie)
            if os.path.isfile(movie_filename) and extension in MOVIE_FILE_TYPES:
                yield movie_name, movie_filename
            elif os.path.isdir(movie_filename):
                yield movie_name, movie_filename

