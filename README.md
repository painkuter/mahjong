# Mahjong

### Installation
    curl https://glide.sh/get | sh
    glide install

### Run
    go run main.go
    
### Links
    /room
    /room/[URL]
    /new-room
    /rooms-list    
    
### Testing
    go test ./...
    
### Statement
* player
    * hand
    * dump
* wall // list of wall's elements
* wind // current wind number (1-east, 2-south, 3-west, 4-north)
* east // east-player number
* step // current player's number
* reserve // 

### Message list
* start game [start]
* stop game [stop]
* update statement [game]
    * skip announce [skip]
    * move tail to discard [discard]
    * announce combination [announce]
        * Chow [chow]
        * Pong [pong]
        * Kong [kong]
            * get tail from reserve
        
* update players list [players]
* message [message] //just text message

### TODO
* ping-pong players
* run room in goroutine
* close ws after game
* makefile
* save history to DB
* use cookies for identification 
* connect to the room by URL
* error's pages
* separate loggers for dev and prod
* save last turn
* remove room from list after player give up