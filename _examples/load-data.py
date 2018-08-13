#!/usr/bin/env python3

import requests

URL = "http://localhost:8080/api/v1/tasks/booru/fetch-tag"

TAGS = [
    "landscape",
    "moon",
    "night",
    "scenic",
    "sky",
    "star",
    "sunset",
    "ruins",
]


def main():
    session = requests.session()
    for tag in TAGS:
        session.post(URL, json={
            "domain": "konachan.net", "tag": tag,
        })


if __name__ == '__main__':
    main()
