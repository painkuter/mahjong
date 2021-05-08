// skip

var textureSkipButton = PIXI.Texture.fromImage('static/images/skip.png');
var textureSkipButtonDown = PIXI.Texture.fromImage('static/images/skip_on.png');
var textureSkipButtonOver = PIXI.Texture.fromImage('static/images/skip.png');

var skipButton = new PIXI.Sprite(textureSkipButton);
skipButton.buttonMode = true;

skipButton.anchor.set(0.5);
skipButton.x = 200;
skipButton.y = 200;

// make the button interactive...
skipButton.interactive = true;
skipButton.buttonMode = true;

skipButton
    // Mouse & touch events are normalized into
    // the pointer* events for handling different
    // button events.
    .on('pointerdown', onButtonDown)
    .on('pointerup', onButtonUp)
    .on('pointerupoutside', onButtonUp)
    .on('pointerover', onButtonOver)
    .on('pointerout', onButtonOut);

// app.stage.addChild(button);

function onButtonDown() {
    this.isdown = true;
    this.texture = textureSkipButtonDown;
    this.alpha = 1;
    sendSkip()
}

function onButtonUp() {
    this.isdown = false;
    if (this.isOver) {
        this.texture = textureSkipButtonOver;
    }
    else {
        this.texture = textureSkipButton;
    }
}

function onButtonOver() {
    this.isOver = true;
    if (this.isdown) {
        return;
    }
    this.texture = textureSkipButtonOver;
}

function onButtonOut() {
    this.isOver = false;
    if (this.isdown) {
        return;
    }
    this.texture = textureSkipButton;
}

function sendSkip(){
    // console.debug("skip")
    const turn = {
        action: "skip",
        meld: "",
        value: []
    };
    const wsMessage = {
        status: "action",
        body: turn
    };
    let req = JSON.stringify(wsMessage)
    console.debug(req)
    conn.send(req);
}

// get

var textureGetButton = PIXI.Texture.fromImage('static/images/get.png');
var textureGetButtonDown = PIXI.Texture.fromImage('static/images/get_on.png');
var textureGetButtonOver = PIXI.Texture.fromImage('static/images/get.png');

var getButton = new PIXI.Sprite(textureGetButton);
getButton.buttonMode = true;

getButton.anchor.set(0.5);
getButton.x = 200;
getButton.y = 250;

// make the button interactive...
getButton.interactive = true;
getButton.buttonMode = true;

getButton
    // Mouse & touch events are normalized into
    // the pointer* events for handling different
    // button events.
    .on('pointerdown', onGetButtonDown)
    .on('pointerup', onGetButtonUp)
    .on('pointerupoutside', onGetButtonUp)
    .on('pointerover', onGetButtonOver)
    .on('pointerout', onGetButtonOut);

function onGetButtonDown() {
    this.isdown = true;
    this.texture = textureGetButtonDown;
    this.alpha = 1;
    sendGetMessage()
}

function onGetButtonUp() {
    this.isdown = false;
    if (this.isOver) {
        this.texture = textureGetButtonOver;
    }
    else {
        this.texture = textureGetButton;
    }
}

function onGetButtonOver() {
    this.isOver = true;
    if (this.isdown) {
        return;
    }
    this.texture = textureGetButtonOver;
}

function onGetButtonOut() {
    this.isOver = false;
    if (this.isdown) {
        return;
    }
    this.texture = textureGetButton;
}

function sendGetMessage(){
    const turn = {
        action: "get_tile",
        meld: "",
        value: []
    };
    const wsMessage = {
        status: "action",
        body: turn
    };
    let req = JSON.stringify(wsMessage)
    console.debug(req)
    conn.send(req);
}