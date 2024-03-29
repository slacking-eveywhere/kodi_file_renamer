#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from pathlib import Path

import re
import unicodedata2
from tvdb import TVDB
from list_movie import MovieList
from list_tvshow import TVShowList, TVShowEpisodesList
from database import TVShowsModel
from exceptions import EpisodePathAlreadyUsedError

DB_NAME = "database.sql"
MOVIE_PATH = "/data-dir/to sort"
MOVIE_PATH_SORTED = "/data-dir/sorted"
TVSHOW_PATH = "/data-dir/series"

UNICODE_NORMALIZE_FORM_VALUE = "NFKC"


def list_movie():
    tvdb = TVDB("9ec9de2268745b801af7c5f21d2a16b8")

    for movie_name, movie_filename, extension, duration in MovieList(MOVIE_PATH):
        movie_propositions = tvdb.search_movies(movie_name)
        for movie_proposition in movie_propositions:
            movie_proposition.directors = tvdb.get_movie_directors_by_movie_id(movie_proposition.id)
            movie_proposition.cast = tvdb.get_movie_cast_by_movie_id(movie_proposition.id)
            movie_proposition.runtime = tvdb.get_movie_runtime_by_id(movie_proposition.id)
        yield movie_name, movie_filename, movie_propositions, extension, duration


def list_tvshow(force_rescan=False):
    tvdb = TVDB("9ec9de2268745b801af7c5f21d2a16b8")
    tvshow_model = TVShowsModel(DB_NAME)
    tvshow_list = TVShowList(TVSHOW_PATH)
    for tvshow_name in tvshow_list:
        if tvshow_model.get_tvshow_by_name(tvshow_name) and force_rescan is False:
            continue
        yield tvshow_name, tvdb.search_tv_shows(
            unicodedata2.normalize(
                UNICODE_NORMALIZE_FORM_VALUE,
                tvshow_name
            )
        )


def list_tvshow_episodes(tvshow_name):
    tvshow_model = TVShowsModel(DB_NAME)
    tvshow = tvshow_model.get_tvshow_by_name(tvshow_name)

    new_episode_path_history = set()

    for episode_path in TVShowEpisodesList(tvshow.path):
        tvshow_episode_details = get_tvshow_episode_details(tvshow.id, episode_path)

        try:
            episode_name = tvshow_episode_details["name"]

            if "/" in episode_name:
                episode_name = episode_name.replace("/", " ")

            new_episode_path = Path(
                episode_path.parent,
                f"""{episode_path.parent.name} S{str(tvshow_episode_details["season_number"]).zfill(2)}E{str(tvshow_episode_details["episode_number"]).zfill(2)} {episode_name}{episode_path.suffix}""")

            tvshow_model.set_or_update_tvshow_episode(
                tvshow_episode_details["id"],
                tvshow_episode_details["name"],
                str(new_episode_path),
                tvshow_episode_details["season_number"],
                tvshow_episode_details["episode_number"],
                tvshow.id
            )

            if new_episode_path in new_episode_path_history:
                raise EpisodePathAlreadyUsedError

            episode_path.rename(new_episode_path)
            new_episode_path_history.add(new_episode_path)

        except KeyError as e:
            print("error on key", tvshow_name, episode_path, e)
            pass
        except EpisodePathAlreadyUsedError:
            print("error on episode name already used", tvshow_name, episode_path, new_episode_path)
            pass


def propose_choice(movies_list):
    for movie_name, movie_filename, movie_propositions, extension, duration in movies_list:
        selected_movie = movie_propositions.choice(int(duration / 60))
        if not selected_movie:
            continue
        if extension:
            create_dir(os.path.join(MOVIE_PATH_SORTED, selected_movie.get_file_name()))
            os.rename(movie_filename, os.path.join(MOVIE_PATH_SORTED, selected_movie.get_file_name(), selected_movie.get_file_name(extension)))
        else:
            os.rename(movie_filename, os.path.join(MOVIE_PATH_SORTED, selected_movie.get_file_name()))
        os.system("clear")


def propose_choice_tv(tvshow_list):
    tvshow_model = TVShowsModel(DB_NAME)
    for tvshow_name, tvshow_propositions in tvshow_list:
        selected_tvshow = tvshow_propositions.choice()
        if not selected_tvshow:
            continue

        current_path = os.path.join(TVSHOW_PATH, tvshow_name)
        new_path = os.path.join(TVSHOW_PATH, selected_tvshow.get_file_name())
        os.rename(current_path, new_path)

        tvshow_model.set_or_update_tvshow(selected_tvshow.id, tvshow_name, new_path)
        os.system("clear")


def get_tvshow_episode_details(tvshow_id, episode_path):
    season, ep_number = None, None

    try:
        ep_number, = re.search("EP([0-9]{1,2})", episode_path.name).groups()
        season = "01"
    except AttributeError:
        pass
    try:
        season, ep_number = re.search("S([0-9]{1,2})E([0-9]{1,2})", episode_path.name, re.IGNORECASE).groups()
    except AttributeError:
        pass
    try:
        season, ep_number = re.search("([0-9]{1,2})x([0-9]{1,2})", episode_path.name, re.IGNORECASE).groups()
    except AttributeError:
        pass
    try:
        season, ep_number = re.search("([0-9]{1,2}) - ([0-9]{1,2})", episode_path.name, re.IGNORECASE).groups()
    except AttributeError:
        pass

    if not season and not ep_number:
        return {}

    return TVDB("9ec9de2268745b801af7c5f21d2a16b8", "en")\
        .get_tvshow_episode_detail_by_id_and_episode_number(str(tvshow_id), season, ep_number)


def create_dir(directory_path):
    try:
        os.mkdir(directory_path)
    except OSError:
        pass


if __name__ == "__main__":
    movies_list = list_movie()
    propose_choice(list(movies_list))

    tvdb_list = list(list_tvshow(True))

    propose_choice_tv(tvdb_list)

    for _tvshow in tvdb_list:
        list_tvshow_episodes(_tvshow[0])

