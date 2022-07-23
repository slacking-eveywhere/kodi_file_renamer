#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import os
from moviepy.editor import VideoFileClip

MOVIE_FILE_TYPES = [".mkv", ".avi", ".mp4", ".iso", ".m2ts", ".img", ".ts"]


class MovieList:
    def __init__(self, movie_path):
        self.movie_path = movie_path

    def __iter__(self):
        ml = sorted(os.listdir(self.movie_path))
        for index, movie in enumerate(ml):
            movie_filename = os.path.join(self.movie_path, movie)
            movie_name, extension = os.path.splitext(movie)
            if os.path.isfile(movie_filename) and extension in MOVIE_FILE_TYPES:
                if extension != ".iso":
                    try:
                        duration = VideoFileClip(movie_filename).duration
                    except UnicodeDecodeError:
                        duration = 0
                else:
                    duration = 0
                yield movie_name, movie_filename, extension, duration
            elif os.path.isdir(movie_filename) and movie_name != "sorted":
                yield movie_name, movie_filename, "", 0
            sys.stdout.write(f"{ index } / { len(ml) }\r")
            sys.stdout.flush()

