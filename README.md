# mahjong

### Installation
    go get github.com/gorilla/websocket

### Run
    github.com/gorilla/websocket
    
### Statement
* player
    * hand
    * dump
* wall // list of wall's elements
* wind // current wind number (1-east, 2-north, 3-west, 4-south)
* step // current player's number
* reserve // 

### Message list
* start game
* stop game
* update statement
* update players list

### TODO
* ping-pong players
* run room in goroutine
* close ws after game