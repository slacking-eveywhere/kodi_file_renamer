#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sqlite3


class Connector:
    def __init__(self, database_name):
        self.database_name = database_name

    def execute(self, sql_command, params=()):
        with sqlite3.connect(self.database_name) as conn:
            cursor = conn.cursor()
            cursor.execute(sql_command, params)
            conn.commit()

    def fetch_one(self, sql_command, param=()):
        with sqlite3.connect(self.database_name) as conn:
            cursor = conn.cursor()
            cursor.execute(sql_command, param)
            return cursor.fetchone()

    def fetch_all(self, sql_command, param=()):
        with sqlite3.connect(self.database_name) as conn:
            cursor = conn.cursor()
            cursor.execute(sql_command, param)
            return cursor.fetchall()


class TVShow:
    def __init__(self, id=-1, name="", path=""):
        self.id = id
        self.name = name
        self.path = path


class TVShowsModel:
    def __init__(self, database_name):
        self.connector = Connector(database_name)
        self.__create_database()

    def __create_database(self):
        self.connector.execute("""CREATE TABLE if not exists tvshow(
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        path TEXT NOT NULL);
        """)
        self.connector.execute("""CREATE TABLE if not exists tvshow_episodes(
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        path TEXT NOT NULL,
        season_number INTEGER NOT NULL,
        episode_number INTEGER NOT NULL,
        tvshow_id INTEGER NOT NULL
        );""")

    def set_or_update_tvshow(self, tvshow_id, tvshow_name, tvshow_path):
        t = (tvshow_id, tvshow_name, tvshow_path)
        return self.connector.execute("""REPLACE INTO tvshow
        VALUES (?, ?, ?);""", t)

    def set_or_update_tvshow_episode(self, episode_id, name, path, season_number, episode_number, tvshow_id):
        t = (episode_id, name, path, season_number, episode_number, tvshow_id)
        return self.connector.execute("""
        REPLACE INTO tvshow_episodes 
        VALUES (?, ?, ?, ?, ?, ?)""", t)

    def get_tvshow_by_name(self, tvshow_name):
        t = (tvshow_name, )
        result = self.connector.fetch_one("""
        SELECT id, name, path
        FROM tvshow
        WHERE name = ?;""", t) or []
        return TVShow(*result)

    def get_tvshow_episodes_by_tvshow_id(self, tvshow_id):
        t = (tvshow_id, )
        result = self.connector.fetch_all("""
        SELECT id, name, path, season_number, episode_number, tvshow_id
        FROM tvshow_episodes
        WHERE tvshow_id = ?""", t)
        return result

    def get_tvshow_episode_by_season_and_episode(self, season_number, episode_number):
        t = (season_number, episode_number)
        result = self.connector.fetch_one("""
        SELECT id, name, path, season_number, episode_number, tvshow_id
        FROM tvshow_episodes
        WHERE season_number = ? AND episode_number = ?""", t)
        return result
