#!/usr/bin/env python3.6

import requests


SERVER_URL = "http://localhost:8080/v1"

MAX_FAVORITES = 100
MAX_POSTS = 100

session = requests.Session()


def test_favorites():
    # create favorites
    for i in range(MAX_FAVORITES):
        response = session.request(
            "POST", f"{SERVER_URL}/favorites", json={"name": f"f:{i:04}"}
        )
        assert response.status_code == 200, "Response is not 200"

    # get all favorites
    response = session.request("GET", f"{SERVER_URL}/favorites")
    favorites = response.json()["favorites"]
    assert response.status_code == 200, "Response is not 200"
    assert (
        len(favorites) == MAX_FAVORITES
    ), f"List size mismatch, found: {len(favorites)}, expect: {MAX_FAVORITES}"

    # delete favorites
    for favorite in favorites:
        response = session.request(
            "DELETE", f"{SERVER_URL}/favorites/{favorite['name']}"
        )
        assert response.status_code == 200, "Response is not 200"

    # assert all favorites deleted
    response = session.request("GET", f"{SERVER_URL}/favorites")
    favorites = response.json()["favorites"]
    assert response.status_code == 200, "Response is not 200"
    assert not favorites, "Not all favorites deleted"


def test_posts():
    # create posts
    for i in range(MAX_POSTS):
        (x, y) = (i // 10, i % 10)
        response = session.request(
            "POST",
            f"{SERVER_URL}/posts",
            json={
                "tags": [f"x:{x:04}", f"y:{y:04}"],
                "namespace": "test",
                "external": True,
                "id": f"k:{i:04}",
                "value": f"v:{i:04}",
            },
        )
        assert response.status_code == 200, "Response is not 200"

    # get all posts
    response = session.request("GET", f"{SERVER_URL}/posts")
    posts = response.json()["posts"]
    assert response.status_code == 200, "Response is not 200"
    assert (
        len(posts) == MAX_POSTS
    ), f"List size mismatch, found: {len(posts)}, expect: {MAX_POSTS}"

    # get all tags
    response = session.request("GET", f"{SERVER_URL}/tags")
    tags = response.json()["tags"]
    for i in range(MAX_POSTS):
        (x, y) = (i // 10, i % 10)
        assert f"x:{x:04}" in tags, f"Tag 'x:{x:04}' not created"
        assert f"y:{y:04}" in tags, f"Tag 'y:{y:04}' not created"

    # delete posts
    for favorite in posts:
        response = session.request(
            "DELETE", f"{SERVER_URL}/posts/{favorite['id']}")
        assert response.status_code == 200, "Response is not 200"

    # assert all tags deleted
    response = session.request("GET", f"{SERVER_URL}/tags")
    tags = response.json()["tags"]
    assert response.status_code == 200, "Response is not 200"
    assert not tags, "Not all tags deleted"

    # assert all posts deleted
    response = session.request("GET", f"{SERVER_URL}/posts")
    posts = response.json()["posts"]
    assert response.status_code == 200, "Response is not 200"
    assert not posts, "Not all posts deleted"


if __name__ == "__main__":
    test_favorites()
    test_posts()
