#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from os import path
from datetime import datetime
import re
import requests
from movie import MoviePropositionsList, MovieProposition
from tvshow import TVShowPropositionsList, TVShowProposition

TVDB_URL = "https://api.themoviedb.org/3/"


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
        print(params)
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

    def get_tvshow_episode_detail_by_id_and_episode_number(self, tvshow_id, season_number, episode_number):
        url = path.join(TVDB_URL, "tv", tvshow_id, "season", season_number, "episode", episode_number)
        req = requests.get(url, self.parameters)

        response = req.json()

        if response.get("success", True) is False:
            return {}

        return response

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
