# Mahjong

#### Installation
    docker-compose up
    go mod tidy

#### Run
    go run main.go
    
#### Links
    0.0.0.0/room
    /room/[URL]
    /new-room
    /rooms-list    
    
#### Testing
    go test ./...
    
#### Statement
* player
    * hand
    * dump
* wall // list of wall's elements
* wind // current wind number (1-east, 2-south, 3-west, 4-north)
* east // east-player number
* step // current player's number
* reserve // 

#### Message list
* start game [start]
* stop game [stop]
* error [error]
* update statement [game]
    * skip announce [skip]
    * move tail to discard [discard]
    * announce combination [announce]
        * Chow [chow]
        * Pong [pong]
        * Kong [kong]
            * get tail from reserve
        * Mahjong [mahjong]
    * announce ready hand
        
* update players list [players]
* message [message] //just text message

#### TODO
* spectate the room
* print YOUR TURN +
* ping-pong players
* run room in goroutine +
* close ws after game
* makefile
* save history to DB
* use cookies for identification 
* connect to the room by URL
* error's pages
* separate logger for dev and prod
* save last turn
* remove room from list after player give up
* monitoring