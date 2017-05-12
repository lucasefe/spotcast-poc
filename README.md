# spotcast

Spotcast allows one spotify desktop client (ideally yours) to manage other spotify desktop clients running on different hosts.

In order to achieve this, you need to install a client binary on each machine running Spotify, and have one server binary running somewhere, available to anyone, and you need to know the url of this one (duh).

## Only one driver!

Who gets to decide what's being played? Simple, the first one to join the channel. The next ones will be the followers. 

## Caveats

Many. 

* If the leader drops from the server, no new leader is elected. Everyone needs to leaver and re-join. This will be fixed soon. 
* No authentication, for now.
* This is a developer thing. It's still pretty early to call it stable, so YMMV. 
* What you can do with each spotify client is pretty limited. You can PLAY, PAUSE, RESUME, and that's it. Playing from a specific position is not supported, nor skipping songs.


## Acknowledgements / Some history

* To Julian Porta, a fellow programmer who played shared/public sessions on Grooveshark every Friday, entertaining us with different kinds of music.  Great times. Now he's not with us anymore, but I remember him everyday, and this is a tribute to him, in some way. 

* I encoutered [this article](http://cgbystrom.com/articles/deconstructing-spotifys-builtin-http-server/) while I was searching for a way to control the spotify client. It gave me all the insight I needed to achieve it. Kudos to the guy who wrote it. 

## Installation 

### Client

You can use homebrew to install it on a mac, otherwise download the binary from the [releases section](https://github.com/lucasefe/spotcast/releases)

    brew tap lucasefe/spotcast
    brew install spotcast-client

To run it, you need to specify the server host:port combination. 

    spotcast-client -remote spotcast-server.herokuapp.com

By default, the clients connects to the `general` channel, but if you want to have a separate broadcast, you can switch to a different one by specifying it:

    spotcast-client -channel heavymetal

Again, if you are the first, you get to decide what to play. 

### Server 

You can run spotcast server on a free heroku instance by clicking [![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Or using homebrew:

    brew tap lucasefe/spotcast
    brew install spotcast-server

## Build // Development

* Copy `cp .env.sample .env`
* Setup environment with `source .env`
* Install dependencies with [gpm](https://github.com/pote/gpm): `gpm install`

All rights reserved to Lucas Florio
