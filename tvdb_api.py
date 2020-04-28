#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from os import path
from datetime import datetime
import requests
from list_movie import MovieList
from movie import MoviePropositionsList, MovieProposition

TVDB_URL = "https://api.themoviedb.org/3/"
MOVIE_PATH = "/Volumes/media/films"
MOVIE_PATH_SORTED = "/Volumes/media/films/sorted"


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
        params = {"query": query}
        params.update(self.parameters)
        req = requests.get(url, params=params)

        return [
            (
                str(result.get("id")),
                result.get("original_name"),
                datetime.strptime(result.get("first_air_date"), "%Y-%m-%d")
            )
            for result in req.json().get("results", [])
        ]

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


def create_dir(dirpath):
    try:
        os.mkdir(dirpath)
    except OSError:
        pass


if __name__ == "__main__":
    movies_list = list_movie()
    propose_choice(list(movies_list))

    # c = a.search_tv_shows("game+of+thrones")

    # id, tv_name, release_date = c[0]
    # print(a.get_tv_shows_by_id(id))

