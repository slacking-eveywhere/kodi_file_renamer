#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from os import path
from datetime import datetime
import requests

TVDB_URL = "https://api.themoviedb.org/3/"


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
            (
                str(result.get("id")),
                result.get("original_title"),
                datetime.strptime(result.get("release_date"), "%Y-%m-%d")
            )
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

    def get_movie_by_id(self, id):
        url = path.join(TVDB_URL, 'movie', id)
        req = requests.get(url, params=self.parameters)

        return req.json()

    def get_tv_shows_by_id(self, id):
        url = path.join(TVDB_URL, 'tv', id)
        req = requests.get(url, params=self.parameters)

        return req.json()


if __name__ == "__main__":
    a = TVDB("9ec9de2268745b801af7c5f21d2a16b8")
    b = a.search_movies("Victoria")
    c = a.search_tv_shows("game+of+thrones")

    id, tv_name, release_date = c[0]
    print(a.get_tv_shows_by_id(id))

