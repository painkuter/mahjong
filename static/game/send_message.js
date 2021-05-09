
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

function moveToDiscard(tile){
    const turn = {
        action: "discard",
        value: [tile]
    };
    const wsMessage = {
        status: "action",
        body: turn
    };
    let req = JSON.stringify(wsMessage)
    console.debug(req)
    conn.send(req);
}