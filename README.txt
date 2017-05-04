

server

  * Listens locally to play/stop/pause requests (http)
  * subscribes to remote server, so remote bot can control it. 

  POST /play/(:trackid)
  POST /stop
  POST /pause
