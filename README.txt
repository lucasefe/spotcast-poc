# spotcast

For now, manage your local spotify.

## CLI


    cli play spotify:album:6eWtdQm0hSlTgpkbw4LaBG
    cli pause
    cli status
    cli search Bon Jovi

## Server

    curl -XGET http://localhost:8081/status
    curl -XGET http://localhost:8081/sessions
    curl -XPOST http://localhost:8081/pause

## TODO

* learn more about the local webser api
