#!/usr/bin/env python3

import requests


POSTS_URL = "http://localhost:8080/v1/posts"
FETCH_URL = "http://localhost:8080/v1/tasks/booru/fetch-file"

QUERY = "eureka_seven | landscape | ruins"


def main():
    session = requests.session()
    r = session.get(POSTS_URL, params={"q": QUERY})
    for row in r.json():
        session.post(FETCH_URL, json={
            "type": row["type"],
            "key": row["key"],
        })


if __name__ == '__main__':
    main()
