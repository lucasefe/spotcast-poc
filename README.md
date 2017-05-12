# spotcast

Spotcast allows one spotify desktop client to drive other spotify desktop clients. 

In order to achieve this, you need to install a client binary on each machine running Spotify, and have one server binary running somewhere, available to anyone. 

## Only one driver!

Who gets to decide what's being played? Simple, the first one to join the server. The next ones will be the followers. 

## Caveats

Many. 

* If the leader drops from the server, no new leader is elected. Everyone needs to leaver and re-join. This will be fixed soon. 
* No authentication, for now. 

## Installation 

### Server 

You can run spotcast server on a free heroku instance by clicking [![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Or using homebrew:

    brew tap lucasefe/spotcast
    brew install spotcast-server

### Client

    brew tap lucasefe/spotcast
    brew install spotcast-client

## Components 

### CLI

    cli play spotify:album:6eWtdQm0hSlTgpkbw4LaBG
    cli pause
    cli status

### Client

    curl -XPOST http://localhost:8080/stop
    curl -XPOST http://localhost:8080/pause
    curl -XPOST http://localhost:8080/play/spotify:artist:08td7MxkoHQkXnWAYD8d6Q

### Server

    curl -XGET http://localhost:8081/status
    curl -XGET http://localhost:8081/sessions
    curl -XPOST http://localhost:8081/pause

## Build // Development

* Copy `cp .env.sample .env`
* Setup environment with `source .env`
* Install dependencies with [gpm](https://github.com/pote/gpm): `gpm install`

All rights reserved to Lucas Florio
