# spotcast

Client/Server application that connects to a centralized server to receive spotify's player requests. 

It can act as leader or follower. 

## CLI

    cli play spotify:album:6eWtdQm0hSlTgpkbw4LaBG
    cli pause
    cli status

## Client

    curl -XPOST http://localhost:8080/stop
    curl -XPOST http://localhost:8080/pause
    curl -XPOST http://localhost:8080/play/spotify:artist:08td7MxkoHQkXnWAYD8d6Q

## Server

    curl -XGET http://localhost:8081/status
    curl -XGET http://localhost:8081/sessions
    curl -XPOST http://localhost:8081/pause

## TODO

* learn more about the local webser api

## Build

* Setup environment with `source .env`
* Install dependencies with [gpm](https://github.com/pote/gpm): `gpm install`
* Optionally, build with [realize](https://github.com/tockins/realize): `realize run`

All rights reserved to Lucas Florio
