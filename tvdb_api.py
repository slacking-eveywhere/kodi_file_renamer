#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from pathlib import Path
import re
from tvdb import TVDB
from list_movie import MovieList
from list_tvshow import TVShowList, TVShowEpisodesList


MOVIE_PATH = "/Volumes/medias/divers/rsync/to sort"
MOVIE_PATH_SORTED = "/Volumes/medias/divers/rsync/sorted"
TVSHOW_PATH = "/Volumes/medias/series"


def list_movie():
    tvdb = TVDB("9ec9de2268745b801af7c5f21d2a16b8")

    for movie_name, movie_filename, extension, duration in MovieList(MOVIE_PATH):
        movie_propositions = tvdb.search_movies(movie_name)
        for movie_proposition in movie_propositions:
            movie_proposition.directors = tvdb.get_movie_directors_by_movie_id(movie_proposition.id)
            movie_proposition.cast = tvdb.get_movie_cast_by_movie_id(movie_proposition.id)
            movie_proposition.runtime = tvdb.get_movie_runtime_by_id(movie_proposition.id)
        yield movie_name, movie_filename, movie_propositions, extension, duration


def list_tvshow():
    tvdb = TVDB("9ec9de2268745b801af7c5f21d2a16b8")
    tvshow_list = TVShowList(TVSHOW_PATH)
    for tvshow_name in tvshow_list:
        yield tvshow_name, tvdb.search_tv_shows(tvshow_name)
        break


def propose_choice(movies_list):
    for movie_name, movie_filename, movie_propositions, extension, duration in movies_list:
        selected_movie = movie_propositions.choice(int(duration / 60))
        if selected_movie:
            if extension:
                create_dir(os.path.join(MOVIE_PATH_SORTED, selected_movie.get_file_name()))
                os.rename(movie_filename, os.path.join(MOVIE_PATH_SORTED, selected_movie.get_file_name(), selected_movie.get_file_name(extension)))
            else:
                os.rename(movie_filename, os.path.join(MOVIE_PATH_SORTED, selected_movie.get_file_name()))
        os.system("clear")


def propose_choice_tv(tvshow_list):
    for tvshow_name, tvshow_propositions in tvshow_list:
        selected_tvshow = tvshow_propositions.choice()
        if selected_tvshow:
            current_path = os.path.join(TVSHOW_PATH, tvshow_name)
            new_path = os.path.join(TVSHOW_PATH, selected_tvshow.get_file_name())
            os.rename(current_path, new_path)

            tvshow_episodes_paths_list = TVShowEpisodesList(new_path)
            set_tvshow_episode_name(tvshow_episodes_paths_list)

        # os.system("clear")


def set_tvshow_episode_name(episodes_paths_list):
    for episode_path in episodes_paths_list:
        season, ep_number = None, None

        try:
            season, ep_number = re.search("S([0-9]{1,2})E([0-9]{1,2})", episode_path.name).groups()
        except AttributeError:
            pass
        try:
            season, ep_number = re.search("s([0-9]{1,2})e([0-9]{1,2})", episode_path.name).groups()
        except AttributeError:
            pass
        try:
            season, ep_number = re.search("([0-9]{1,2})x([0-9]{1,2})", episode_path.name).groups()
        except AttributeError:
            pass

        parent_path = episode_path.parent

        if not season and not ep_number:
            continue

            # a = TVDB("9ec9de2268745b801af7c5f21d2a16b8", "en").get_tvshow_episode_detail_by_id_and_episode_number(selected_tvshow.id, season, ep_number)
            # print(a.get("name"))
        new_episode_path = Path(parent_path, f"{parent_path.name} S{season}E{ep_number}{episode_path.suffix}")
        episode_path.rename(new_episode_path)


def create_dir(directory_path):
    try:
        os.mkdir(directory_path)
    except OSError:
        pass


if __name__ == "__main__":
    # movies_list = list_movie()
    # propose_choice(list(movies_list))

    tvdb_list = list_tvshow()
    propose_choice_tv(list(tvdb_list))
