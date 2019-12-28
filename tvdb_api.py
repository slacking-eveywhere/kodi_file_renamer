#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
from pprint import pprint

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
        result = requests.get(url, params=params)
        return result.json()


if __name__ == "__main__":
    a = TVDB("9ec9de2268745b801af7c5f21d2a16b8")
    b = a.search_movie("VERY BAD TRIP")
    pprint(b)
