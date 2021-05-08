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
        var tile = tileStack[player.discard[i]];

        tile.sprite.width = 40;
        tile.sprite.height = 56;
        tile.sprite.x = (tile.sprite.width + 2) * col;
        tile.sprite.y = (tile.sprite.height + 2) * row;
        discard.addChild(tile.sprite);

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

function load() {
    for (const index in tiles_data.mahjong_tiles) {
        const value = tiles_data.mahjong_tiles[index];
        for (let i = 1; i <= 4; i++) {
            tileStack[index + "_" + i] = createTile(PIXI.loader.resources.mahjongTiles.texture, new PIXI.Rectangle(value.x, value.y, value.width, value.height), index + "_" + i);
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