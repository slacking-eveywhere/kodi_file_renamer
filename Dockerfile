FROM python:3.9-slim

COPY requirements.txt requirements.txt

RUN apt-get update ; \
    apt-get upgrade ; \
    apt-get install -y \
        git \
        ffmpeg

RUN python -m pip install -r requirements.txt
