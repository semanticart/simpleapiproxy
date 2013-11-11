simpleapiproxy
============

A small proxy to allow consuming apis via javascript without exposing the api keys.

As an example, I'm using this with last.fm for a testbed.

example config:

    URL_ROOT=http://ws.audioscrobbler.com/2.0/ URL_SUFFIX=api_key=XXXXXXXXXXXXX PORT=80

This would serve http://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=violencenow&format=json&api_key=XXXXXXXXXXXXX when you request http://example.com/?method=user.getrecenttracks&user=violencenow&format=json

deploying
============

Just folow the instructions for adding the heroku buildpack and you should be good to go: https://github.com/kr/heroku-buildpack-go

Procfile included.
