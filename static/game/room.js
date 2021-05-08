// alert('from room.js')


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
    console.debug("here");
    let me = getMe();
    app.stage.addChild(me);
    //Render the stage
    app.renderer.render(app.stage);
}

function getMe() {
    let me = new PIXI.Container();
    console.debug("hand "+gameData.body.players["100"].hand);
    for (var i in gameData.body.players["100"].hand) {
        var tile = tileStack[gameData.body.players["100"].hand[i]];
        tile.sprite.interactive = true;
        tile.sprite.buttonMode = true;
        tile.sprite.on("pointerdown", onSelect);
        console.debug("tile "+tile)
        myStack.push(tile);
    }
    myStack.sort(handSort);

    let discard = getDiscard(gameData.body.players["100"]); //
    console.debug(discard);
    me.addChild(discard);
    console.debug("myStack")
    console.debug(myStack)
    for (var i in myStack) {
        myStack[i].sprite.x = (myStack[i].sprite.width + 2) * i;
        myStack[i].sprite.y = discard.height + 100;
        me.addChild(myStack[i].sprite);
    }

    discard.position.x = me.width / 2 - discard.width / 2;

    me.position.x = app.renderer.view.width / 2 - me.width / 2;
    me.position.y = app.renderer.view.height - (me.height + 10);

    return me;
}

function getDiscard(player) {
    console.debug("player"+player);
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
        console.debug(value)
        for (let i = 1; i <= 4; i++) {
            tileStack[index + "_" + i] = createTile(PIXI.loader.resources.mahjongTiles.texture, new PIXI.Rectangle(value.x, value.y, value.width, value.height), index + "_" + i);
        }
    }
}

function handSort(a, b) {
    // console.log(a.name);
    // console.log(b.name);
    if (a.name > b.name) {
        return 1;
    }
    if (a.name < b.name) {
        return -1;
    }
}