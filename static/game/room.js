function createTile(texture, rect, name) {

    let sprite = new PIXI.Sprite(new PIXI.Texture(texture, rect, name));
    sprite.width = 50;
    sprite.height = 70;
    return {
        "name": name,
        "sprite": sprite
    };
}

function setup() {
    load();
    let me = getMe();
    containers[100] = me
    app.stage.addChild(me);
    app.stage.addChild(skipButton);
    app.stage.addChild(getButton);

    //Render the stage
    app.renderer.render(app.stage);
}

function getMe() {
    let me = new PIXI.Container();
    for (var i in getMyData().hand) {
        var tile = tileStack[getMyData().hand[i]];
        tile.sprite.interactive = true;
        tile.sprite.buttonMode = true;
        tile.sprite.on("pointerdown", onSelect);
        myStack.push(tile);
    }
    myStack.sort(handSort);

    let discard = getDiscard(getMyData()); //
    me.addChild(discard);
    for (var i in myStack) {
        myStack[i].sprite.x = (myStack[i].sprite.width + 2) * i;
        if (getMyData().discard.length > 0) {
            myStack[i].sprite.y = discard.height + 100;
        }
        me.addChild(myStack[i].sprite);
    }

    if (getMyData().discard.length > 0) {
        discard.position.x = me.width / 2 - discard.width / 2;
    }
    me.position.x = app.renderer.view.width / 2 - me.width / 2;
    me.position.y = app.renderer.view.height - (me.height + 10);

    return me;
}

function getMyData() {
    return gameData.body.players["100"]
}

function getDiscard(player) {
    let discard = new PIXI.Container();
    var col = 0;
    var row = 0;
    for (var i in player.discard) {
        var discardTile = tileStack[player.discard[i]];

        discardTile.sprite.width = OPEN_TILE_WIDTH;
        discardTile.sprite.height = OPEN_TILE_HEIGHT;
        discardTile.sprite.x = (discardTile.sprite.width + 2) * col;
        discardTile.sprite.y = (discardTile.sprite.height + 2) * row;
        discard.addChild(discardTile.sprite);

        if (i % 5 == 0 && i > 0) {
            col = 0;
            row++;
        } else {
            col++;
        }
    }
    return discard;
}

function onSelect() {
    if (!this.isSelected) {
        for (const index in myStack) {
            const tile = myStack[index];
            if (tile.sprite.isSelected) {
                tile.sprite.y += 15;
                tile.sprite.isSelected = false;
                break;
            }
        }
        this.isSelected = true;
        this.y -= 15;
    }
}

function load() { // грузит картинки тайлов в память
    for (const index in tiles_data.mahjong_tiles) {
        const value = tiles_data.mahjong_tiles[index];
        for (let i = 1; i <= 4; i++) {
            tileStack[index + "_" + i] = createTile(PIXI.Loader.shared.resources.mahjongTiles.texture, new PIXI.Rectangle(value.x, value.y, value.width, value.height), index + "_" + i);
        }
    }
}

function handSort(a, b) {
    if (a.name > b.name) {
        return 1;
    }
    if (a.name < b.name) {
        return -1;
    }
}

function processDblClick(point) { // dblclick moves tile to discard
    for (const i in myStack) {
        if (myStack[i].sprite.containsPoint(point)) {
            // TODO disable moving tiles to discard
            moveToDiscard(myStack[i].name)
        }
    }
}


function onAction(message) {
    switch (message.body.action) {
        case "announce":
            onAnnounce(message.body, message.body.player);
            break
        case "discard":
            onDiscard(message.body, message.body.player);
            break
        case "get_tile":
            onGetTile(message.body, message.body.player);
            break
        case "skip":
            onSkip(message.body, message.body.player);
            break
    }
}

function onAnnounce(action, playerID) {
    let announce = {
        meld: action.meld,
        value: []
    };
    for (let i in action.value) {
        announce.value.push(action.value[i]);
    }
    if (gameData.body.players[playerID].open == null) {
        gameData.body.players[playerID].open = [];
    }
    gameData.body.players[playerID].open.push(announce);

    drawOpenTiles(playerID)
}

function drawOpenTiles(playerID) {
    let col = 0;
    let offset = 0; // сдвиг относительно предыдущей комбинации


    // console.debug(playerID)
    for (let i in gameData.body.players[playerID].open) {
        for (let j in gameData.body.players[playerID].open[i].value) {

            let tile = tileStack[gameData.body.players[playerID].open[i].value[j]];
            if (tile == null) {
                console.error("Tile not found " + gameData.body.players[playerID].open[i].value[j]);
            }
            tile.sprite.width = OPEN_TILE_WIDTH;
            tile.sprite.height = OPEN_TILE_HEIGHT;
            tile.sprite.x = /*openPositionByPlayerID[playerID].x +*/ (tile.sprite.width + 2) * col + offset;
            tile.sprite.y = /*openPositionByPlayerID[playerID].y +*/ (tile.sprite.height + 2);
            openTiles[playerID].addChild(tile.sprite);
            col++;
        }
        offset += 30
    }
}

function onDiscard(action, playerID) {
    console.debug("Discard" + action.value);

    let smallDiscardWidth = OPEN_TILE_WIDTH * 0.8
    let smallDiscardHeight = OPEN_TILE_HEIGHT * 0.8

    let discardSize = discardContainer.children.length

    // console.debug("size" + discardSize)
    // уменьшаем предыдущий тайл в дискарде
    if (discardSize > 0) {
        let lastTile = discardContainer.getChildAt(discardSize - 1)
        lastTile.width = smallDiscardWidth
        lastTile.height = smallDiscardHeight
    }

console.debug("x_offset"+ (discardSize % DISCARD_CALS))
console.debug("y_offset"+ ~~(discardSize / DISCARD_CALS))
    let discard = tileStack[action.value]
    discard.sprite.x = (smallDiscardWidth + 2) * (discardSize % DISCARD_CALS)
    discard.sprite.y = (smallDiscardHeight + 2) * (~~(discardSize / DISCARD_CALS))
    discard.sprite.width = OPEN_TILE_WIDTH * 1.5
    discard.sprite.height = OPEN_TILE_HEIGHT * 1.5
    discardContainer.addChild(discard.sprite)
    app.renderer.render(app.stage);
}

function onGetTile(action, playerID) {

}

function onSkip(action, playerID) {
}

/*
const openPositionByPlayerID = {
    "1": {
        x: 100,
        y: 400
    },
    "2": {
        x: 600,
        y: 100
    },
    "3": {
        x: 500,
        y: 0
    },
    "100": {
        x: 300,
        y: 300
    },
}
*/

