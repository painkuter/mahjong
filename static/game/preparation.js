// логирование
window.logger = document.getElementById("log");
console.debug = function (msg) {
    window.logger.innerHTML += msg + '\n';
    console.log(msg);
}
// проверка поддержки WebSocket
if (!window["WebSocket"]) {
    alert("WebSocket problem")
}

// проверка доступности WebGL
let type = "WebGL";
if (!PIXI.utils.isWebGLSupported()) {
    type = "canvas"
    alert('canvas');
}

const WIDTH = 1600
const HEIGHT = 900
const BORDER_PADDING_X = 80
const BORDER_PADDING_Y = 110
const OPEN_TILE_WIDTH = 40
const OPEN_TILE_HEIGHT = 56
const DISCARD_CALS = 4

// создание игровой зоны PIXI
window.PIXI = PIXI;
PIXI.utils.sayHello(type);
const app = new PIXI.Application({width: WIDTH, height: HEIGHT});
app.renderer.backgroundColor = 0xeeffee;
//Add the canvas that Pixi automatically created for you to the HTML document
document.getElementById("game").appendChild(app.view);

conn = wsConn();

conn.onopen = function () {
    console.debug("Connected");
};
conn.onclose = function (evt) {
    console.debug("Connection closed");
};

conn.onmessage = function (evt) {
    //parse data to JSON
    let message = JSON.parse(evt.data);
    processMsg(message);
}

// глобальные объекты
let tileStack = {}; // все тайлы
let myStack = []; // мой стек
let gameData; // общее состояние игры, известное текущему игроку

let playerIDs = ["1", "2", "3", "100"]


let openTiles = {
    "100": new PIXI.Sprite(PIXI.Texture.WHITE),
    "1": new PIXI.Sprite(PIXI.Texture.WHITE),
    "2": new PIXI.Sprite(PIXI.Texture.WHITE),
    "3": new PIXI.Sprite(PIXI.Texture.WHITE)
};

let containers = buildPlayerContainers(playerIDs)
let discardContainer = new PIXI.Container()
discardContainer.x = 300
discardContainer.y = 200

const discardLabel = new PIXI.Text("Discard");
discardLabel.x = 50;
discardLabel.y = 100;
discardContainer.addChild(discardLabel)

app.stage.addChild(discardContainer)

containers["100"].x = 0
containers["100"].y = HEIGHT * 0.66

containers["1"].x = BORDER_PADDING_X
containers["1"].y = BORDER_PADDING_Y

containers["2"].x = (WIDTH - BORDER_PADDING_X - (OPEN_TILE_WIDTH + 2) * 3) / 2
containers["2"].y = 0

containers["3"].x = WIDTH - BORDER_PADDING_X - (OPEN_TILE_WIDTH + 2) * 3
containers["3"].y = BORDER_PADDING_Y

setOpenTiles(containers)
setContainers(containers)

app.renderer.render(app.stage);

// открытые комбинации по игрокам  (playerID -> [[1_1_1,1_1_2,1_1_3],[2_1_1,2_1_2,2_1_3]]

function wsConn() { // новое подключение
    const name = "pName"; // TODO запрашивать имя

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
        case "action":
            onAction(message);
            break
    }
    app.renderer.render(app.stage);
}

function onStart() {
    console.debug("start");
}

function onStop(message) {
    console.debug("player " + message.body + " gave up");
}

function onGame(message) {
    gameData = message;

    PIXI.Loader.shared
        .add("mahjongTiles", "/static/images/tiles-mahjong.jpg")
        .load(setup);

    setPlayerName(containers)
}


const interactionManager = app.renderer.plugins.interaction;
app.view.addEventListener("dblclick", (e) => {
    const global = new PIXI.Point();
    interactionManager.mapPositionToPoint(global, e.clientX, e.clientY);
    processDblClick(global)
});

function buildPlayerContainers(playerIDs) {
    let containers = {}
    for (let k in playerIDs) {
        containers[playerIDs[k]] = new PIXI.Container()
    }
    return containers
}

function setPlayerName(containers) {
    for (let k in containers) {

        const playerName = new PIXI.Text(gameData.body.players[k].name);
        playerName.x = 50;
        playerName.y = 100;
        containers[k].addChild(playerName)

        console.debug(containers[k]);
    }
}

function setOpenTiles(containers) {
    for (let k in containers) {
        containers[k].addChild(openTiles[k])
    }
}

function setContainers(containers) {
    for (let k in containers) {
        app.stage.addChild(containers[k])
    }
}