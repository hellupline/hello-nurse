# Hello-Nurse

## What is this

Hello-Nurse is a toy project I created to query using set operations on
tags.


## How to use it

create objects using httpie:

sh```
http :8080/v1/posts \
    namespace="test" \
    id="id" \
    value="value" \
    tags="tag1" \
    tags="tag2"
```

query using httpie:


sh```
http :8080/v1/posts q=='(tag1 & tag2)'
http :8080/v1/posts q=='(tag1 | tag2)'
http :8080/v1/posts q=='(tag1 ^ tag2)'
```

accepted operations are intersection (&), union (|) and difference (^)


2 examples of how to push data can be found, just `make run` inside the
push folder, or run the python3.6 script `tests.py`


## TODO
* Better handling of operations, and parser errors
* move http to main.go
* rename push to booru-push
* create booru-serve
* Single Page Application
  * List of session picked files
  * List of Favorites
  * List of Posts per query
