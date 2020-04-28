#!/usr/bin/env python3
# -*- coding: utf-8 -*-


from PyInquirer import prompt, print_json, Separator


class MoviePropositionsList:
    def __init__(self, movie_name, propositions):
        self.movie_name = movie_name
        self.selected_choice = -1
        self.__propositions = propositions

    def __iter__(self):
        for proposition in self.__propositions:
            yield proposition

    def has_been_selected(self):
        return self.selected_choice > -1

    def display_details(self, duration):
        for index, proposition in enumerate(self.__propositions):
            runtime = proposition.runtime if proposition.runtime else 0
            duration = duration if duration else 0
            if duration == 0 or abs(duration - runtime) < 1500:
                yield {
                    "value": index,
                    "name": proposition.display_details()
                }
                yield Separator()

    def choice(self, duration=0):
        if self.__propositions:
            choices = list(self.display_details(duration))
            if not choices:
                choices = list(self.display_details(0))
            choices.append({"value": -1, "name": "None"})

            questions = [
                {
                    'type': 'list',
                    'name': 'propositions',
                    'message': f"What is the correct proposition for '{ self.movie_name }' ({ duration } mn)?",
                    'choices': choices
                }
            ]

            try:
                answers = prompt(questions)
                if answers["propositions"] > -1:
                    return self.__propositions[answers["propositions"]]
                else:
                    return None
            except (ValueError, IndexError):
                return None
        else:
            print(f"No proposition found for '{ self.movie_name }'")


class MovieProposition:
    def __init__(self, _id, title, original_title, language, overview, release_date=None, directors=[], cast=[], runtime=None):
        self.id = _id
        self.title = title
        self.original_title = original_title
        self.language = language
        self.overview = overview
        self.release_date = release_date
        self.directors = directors
        self.cast = cast
        self.runtime = runtime

    def display_details(self):
        return f"""id              :   { self.id }
   title           :   { self.title }
   original title  :   { self.original_title }
   language        :   { self.language}
   release date    :   { self.release_date }
   runtime         :   { self.runtime }
   directors       :   { ", ".join(self.directors) }
   cast            :   { ", ".join(self.cast) }
   overview        :   { self.overview }"""

    def get_file_name(self, extension=""):
        original_title_corrected = self.original_title.replace(":", "-")
        original_title_corrected = original_title_corrected.replace("/", " ")
        return f"{ original_title_corrected } ({ self.release_date.year }){ extension }"
