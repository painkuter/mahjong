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

#### WS requests
    game
    
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
* starting statement [statement]
* update statement [action]
    * gesture [gesture] просясящая рука
    * skip announce [skip]
    * move tail to discard [discard]
    * announce combination [announce]
        * Chow [chow]
        * Pong [pong]
        * Kong [kong]
            * get tail from reserve
        * Mahjong [mahjong]
    * announce ready hand [ready]
        
* update players list [players]
* message [message] //just text message

#### TODO
* Design
  * драконы
* FE
  * наполнение gameData с бэкенда
  * таймер
  * комбинации в открытую (дискард общий, в комбинации отдельно перед каждым)
  * остальные игроки
* BE
  * Name generator 
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
  * create DB for statements, games and players saving

Подключение 
0. /room -> room_name
0. /ws?room=room_name
0. После подключения 4х игроков начинается игра - сообщение "start"-type
0. Пересылка всего состояния - "game"-type

### Описание
* Отношение сторон для браузера фиксированное - 16\9
* Дискард общий, сортируется в порядке сброса
* Последний тайл выделяется в ход, когда он был сброшен
* Если тайл можно взять, выделять еще сильнее