#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from os import path
from datetime import datetime
import requests
from list_movie import MovieList

TVDB_URL = "https://api.themoviedb.org/3/"
MOVIE_PATH = "/Volumes/media/films"


class TVDB:
    def __init__(self, apikey, language="fr"):
        self.parameters = {
            "api_key": apikey,
            "language": language
        }

    def search_movies(self, query):
        url = path.join(TVDB_URL, 'search', 'movie')
        params = {"query": query}
        params.update(self.parameters)
        req = requests.get(url, params=params)

        return [
            {
                "id": str(result.get("id")),
                "title": result.get("title"),
                "original_title": result.get("original_title"),
                "language": result.get("original_language"),
                "overview": result.get("overview"),
                "release_date": parse_release_date(result.get("release_date"))
            }
            for result in req.json().get("results", [])
        ]

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
        return [members.get("name") for members in movie.get("cast")]

    def get_movie_directors_by_movie_id(self, id):
        movie = self.get_movie_credits_by_id(id)
        return [member.get("name") for member in movie.get("crew") if member.get("job") == "Director"]

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


if __name__ == "__main__":
    tvdb = TVDB("9ec9de2268745b801af7c5f21d2a16b8")

    for movie_name, movie_filename in MovieList(MOVIE_PATH):
        movie_results = tvdb.search_movies(movie_name.replace(" ", "+"))
        for movie_result in movie_results:
            movie_result["directors"] = tvdb.get_movie_directors_by_movie_id(movie_result["id"])
            movie_result["cast"] = tvdb.get_movie_cast_by_movie_id(movie_result["id"])
            print(movie_result)
        break
    # c = a.search_tv_shows("game+of+thrones")

    # id, tv_name, release_date = c[0]
    # print(a.get_tv_shows_by_id(id))

