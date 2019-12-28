#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from datetime import datetime
import requests

TVDB_URL = "https://api.themoviedb.org/3/"


class TVDB:
    def __init__(self, apikey, language="fr"):
        self.parameters = {
            "api_key": apikey,
            "language": language
        }

    def search_movie(self, query):
        url = TVDB_URL + 'search/movie'
        params = {"query": query}
        params.update(self.parameters)
        req = requests.get(url, params=params)
        return [
            (
                result.get("original_title"),
                datetime.strptime(result.get("release_date"), "%Y-%m-%d")
            )
            for result in req.json().get("results", [])
        ]


if __name__ == "__main__":
    a = TVDB("9ec9de2268745b801af7c5f21d2a16b8")
    b = a.search_movie("Victoria")
    print(b)