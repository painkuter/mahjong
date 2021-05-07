let type = "WebGL";
if(!PIXI.utils.isWebGLSupported()){
    type = "canvas"
    alert('canvas');

}

let tileStack = {}; // сброс
let myStack = [];
PIXI.utils.sayHello(type);
let app = new PIXI.Application({width: 1200, height: 900});
app.renderer.backgroundColor = 0xffffff;
//Add the canvas that Pixi automatically created for you to the HTML document
document.body.appendChild(app.view);

PIXI.loader
    .add("mahjongTiles", "/static/images/tiles-mahjong.jpg")
    .load(setup);
