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

// создание игровой зоны PIXI
window.PIXI = PIXI;
PIXI.utils.sayHello(type);
const app = new PIXI.Application({width: 1600, height: 500});
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
    var message = JSON.parse(evt.data);
    console.debug("Parsed");
    processMsg(message);
}

// глобальные объекты
let tileStack = {}; // все тайлы
let myStack = []; // мой стек
let gameData; // общее состояние игры, известное текущему игроку
let openTiles = {
    "100": new PIXI.Container(),
    "1": new PIXI.Container(),
    "2": new PIXI.Container(),
    "3": new PIXI.Container()
};

app.stage.addChild(openTiles["100"])
app.stage.addChild(openTiles["1"])
app.stage.addChild(openTiles["2"])
app.stage.addChild(openTiles["3"])
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
}


const interactionManager = app.renderer.plugins.interaction;
app.view.addEventListener("dblclick", (e) => {
    const global = new PIXI.Point();
    interactionManager.mapPositionToPoint(global, e.clientX, e.clientY);
    processDblClick(global)
});
