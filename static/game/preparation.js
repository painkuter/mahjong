console.debug = function(msg) {
    window.logger.innerHTML += msg+'\n';
    console.log(msg);
}

let type = "WebGL";
if (!PIXI.utils.isWebGLSupported()) {
    type = "canvas"
    alert('canvas');

}

if (!window["WebSocket"]) {
    alert("WebSocket problem")
}

window.logger = document.getElementById("log");

conn = wsConn();

conn.onopen = function () {
    console.debug("Connected");
};
conn.onclose = function (evt) {
    console.debug("Connection closed");
};

conn.onmessage = function (evt) {
    //parse data to JSON
    var message = JSON.parse(evt.data);
    console.debug("Parsed");
    processMsg(message);
}

let tileStack = {}; // сброс
let myStack = [];
PIXI.utils.sayHello(type);
let app = new PIXI.Application({width: 1600, height: 500});
app.renderer.backgroundColor = 0xeeffee;
//Add the canvas that Pixi automatically created for you to the HTML document
document.body.appendChild(app.view);

function wsConn() {
    const name = "pName";

    let roomURL = window.location.href.substring(window.location.href.indexOf(host) + host.length + 5); //rewrite it
    console.debug(roomURL);
    return new WebSocket("ws://" + host + "/ws?room=" + roomURL + "&name=" + name);
}

function processMsg(message) {
    switch (message.status) {
        case "start":
            onStart();
            break;
        case "stop":
            onStop(message);
            break;
        case "game":
            onGame(message);
            break;
    }
}

function onStart() {
    console.debug("start");
}

function onStop(message) {
    console.debug( "player " + message.body + " gave up");
}

function log(text) {
    window.logger.innerHTML += text+'\n';
    console.log(text);
}

function onGame(message){
    console.debug(message);
    console.debug(gameData);
    gameData = message;
    // log("message"+JSON.stringify(message));
    // log("gamedata"+JSON.stringify(gameData));
    PIXI.loader
        .add("mahjongTiles", "/static/images/tiles-mahjong.jpg")
        .load(setup);
    console.debug("game started");
}

