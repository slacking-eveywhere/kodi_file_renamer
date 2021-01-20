#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from os import path
from datetime import datetime
import re
import pathlib
import requests
from list_movie import MovieList
from list_tvshow import TVShowList
from movie import MoviePropositionsList, MovieProposition
from tvshow import TVShowPropositionsList, TVShowProposition

TVDB_URL = "https://api.themoviedb.org/3/"
MOVIE_PATH = "/Volumes/medias/divers/rsync/to sort"
MOVIE_PATH_SORTED = "/Volumes/medias/divers/rsync/sorted"
TVSHOW_PATH = "/Volumes/medias/series"


class TVDB:
    def __init__(self, apikey, language="fr"):
        self.parameters = {
            "api_key": apikey,
            "language": language
        }

    def search_movies(self, query):
        url = path.join(TVDB_URL, 'search', 'movie')
        formatted_query = query.replace(" ", "+")
        params = {"query": formatted_query}
        params.update(self.parameters)
        req = requests.get(url, params=params)

        return MoviePropositionsList(query, [
            MovieProposition(**{
                "_id": str(result.get("id")),
                "title": result.get("title"),
                "original_title": result.get("original_title"),
                "language": result.get("original_language"),
                "overview": result.get("overview"),
                "release_date": parse_release_date(result.get("release_date"))
            })
            for result in req.json().get("results", [])
        ])

    def search_tv_shows(self, query):
        url = path.join(TVDB_URL, 'search', 'tv')

        if re.search("\(([0-9]*?)\)", query):
            params = {
                "query": query[:-6],
                "first_air_date_year": query[-5:-1]
            }
        else:
            params = {"query": query}
        params.update(self.parameters)
        req = requests.get(url, params=params)

        return TVShowPropositionsList(query, [
            TVShowProposition(**{
                "_id": str(result.get("id")),
                "title": result.get("name"),
                "original_title": result.get("original_name"),
                "language": result.get("original_language"),
                "overview": result.get("overview"),
                "release_date": parse_release_date(result.get("first_air_date"))
            })
            for result in req.json().get("results", [])
        ])

    def get_movie_overview_by_id(self, id):
        movie = self.get_movie_by_id(id)
        return movie.get("overview")

    def get_movie_cast_by_movie_id(self, id):
        movie = self.get_movie_credits_by_id(id)
        return [members.get("name") for members in movie.get("cast", [])]

    def get_movie_directors_by_movie_id(self, id):
        movie = self.get_movie_credits_by_id(id)
        return [member.get("name") for member in movie.get("crew", []) if member.get("job") == "Director"]

    def get_movie_runtime_by_id(self, id):
        movie = self.get_movie_by_id(id)
        return movie.get("runtime")

    def get_movie_by_id(self, id):
        url = path.join(TVDB_URL, 'movie', id)
        req = requests.get(url, params=self.parameters)

        return req.json()

    def get_movie_credits_by_id(self, id):
        url = path.join(TVDB_URL, 'movie', id, 'credits')
        req = requests.get(url, params=self.parameters)

        return req.json()

    def get_tv_shows_by_id(self, id):
        url = path.join(TVDB_URL, 'tv', id)
        req = requests.get(url, params=self.parameters)

        return req.json()


def parse_release_date(release_date):
    try:
        return datetime.strptime(release_date, "%Y-%m-%d")
    except (ValueError, TypeError):
        return "unknown"


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


def propose_choice(movies_list):
    for movie_name, movie_filename, movie_propositions, extension, duration in movies_list:
        selected_movie = movie_propositions.choice(int(duration / 60))
        if selected_movie:
            print(movie_filename, os.path.join(MOVIE_PATH_SORTED, selected_movie.get_file_name(extension)))
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
            if current_path != new_path:
                os.rename(current_path, new_path)

            old_paths = []
            new_paths = []

            for episode in sorted(os.listdir(new_path)):
                if not episode.startswith("."):
                    episode_path = os.path.join(new_path, episode)
                    ext = pathlib.Path(episode_path).suffix

                    season, ep_number = None, None

                    try:
                        season, ep_number = re.search("S([0-9]{1,2})E([0-9]{1,2})", episode).groups()
                    except AttributeError:
                        pass
                    try:
                        season, ep_number = re.search("s([0-9]{1,2})e([0-9]{1,2})", episode).groups()
                    except AttributeError:
                        pass
                    try:
                        season, ep_number = re.search("([0-9]{1,2})x([0-9]{1,2})", episode).groups()
                    except AttributeError:
                        pass

                    if season and ep_number:
                        new_name = os.path.join(new_path, f"{ selected_tvshow.get_file_name() } S{ season }E{ ep_number }{ ext }")
                        # old_paths.append(episode_path)
                        # new_paths.append(new_name)
                        os.rename(episode_path, new_name)
            # for index, pa in enumerate(old_paths):
            #     print(pa, new_paths[index])
            # input("bla")

        os.system("clear")


def create_dir(dirpath):
    try:
        os.mkdir(dirpath)
    except OSError:
        pass


if __name__ == "__main__":
    # movies_list = list_movie()
    # propose_choice(list(movies_list))

    tvdb_list = list_tvshow()
    propose_choice_tv(list(tvdb_list))

    # id, tv_name, release_date = c[0]
    # print(a.get_tv_shows_by_id(id))

