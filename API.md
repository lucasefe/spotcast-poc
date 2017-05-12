# Spotcast API

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

    curl -XGET http://localhost:8081/channel/:name
