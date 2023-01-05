window.onload = initWebSocket
let globalWs;
let loggedInUserEmail;
let loggedInUserName;
const WsCommands = {
    MESSAGE: "MESSAGE",
    STOCK_TICKER: "STOCK_TICKER",
    ROOM_CREATE: "ROOM_CREATE",
    ROOM_READ: "ROOM_READ",
    USERS: "USERS",
    BAD_SESSION: "BAD_SESSION",
    HISTORY: "HISTORY",
};

const currentRoom = "Public";

function initWebSocket() {
    globalWs = new WebSocket("ws://localhost:6001")
    globalWs.onopen = function () {
        console.log("ws connection opened")
        webSocketCommand(WsCommands.USERS)
        webSocketCommand(WsCommands.ROOM_READ) //get default room list
        webSocketCommand(WsCommands.HISTORY) //get SELECTED ROOM MESSAGES

    }
    globalWs.onmessage = webSocketMMessageHandler
}

function webSocketCommand(command, newMsg) {
    const urlParams = new URLSearchParams(window.location.search);
    const sessionId = urlParams.get('sessionid');
    /*
        //message schema
        const WsMessage = {
            data: "",
            command: "",
            timestamp: "",
            room: "",
            sessionKey: "",
            stockCode:""
        }
    */

    let WsMessage = {
        sessionKey: sessionId
    }

    switch (command) {
        case WsCommands.USERS: //fetch online users
            WsMessage.command = WsCommands.USERS
            break;
        case WsCommands.ROOM_READ: //fetch online users
            WsMessage.command = WsCommands.ROOM_READ
            break;
        case WsCommands.HISTORY: //fetch messages for selected room
            WsMessage.command = WsCommands.HISTORY
            WsMessage.room = currentRoom
            break;

        case WsCommands.MESSAGE: //send a new chat message for the current room
            WsMessage.command = WsCommands.MESSAGE
            WsMessage.room = currentRoom
            WsMessage.data = JSON.stringify({
                name: loggedInUserName,
                email: loggedInUserEmail,
                data: newMsg,
                timestamp: (new Date()).toJSON()
            })
            break;
        default:
    }
    globalWs.send(JSON.stringify(WsMessage))

}

function webSocketMMessageHandler(event) {
    let evData = JSON.parse(event.data)
    switch (evData.command) {
        case WsCommands.USERS:

            //get returned users list
            let users = JSON.parse(evData.data)
            insertOnlineUsers(users)
            break;
        case WsCommands.ROOM_READ:

            let rooms = JSON.parse(evData.data)
            insertRoomList(rooms)
            break;
        case WsCommands.HISTORY:

            let roomMsg = JSON.parse(evData.data)
            insertMessageHistory(roomMsg)
            break;
        case WsCommands.MESSAGE:

            if (evData.room !== currentRoom) break;
            let newMsg = JSON.parse(evData.data)
            insertNewMessage(newMsg)
            break;
        case WsCommands.BAD_SESSION:

            //bad sessionid
            alert("session not valid, please login")
            location.assign("index.html")
            break;
    }

}

function insertOnlineUsers(usersArray) {
    const urlParams = new URLSearchParams(window.location.search);
    const sessionId = urlParams.get('sessionid');

    const userDOMList = document.getElementById("chat-users-online");
    userDOMList.innerHTML = ""
    let sessionKey;
    for (let i = 0; i < usersArray.length; i++) {
        userDOMList.innerHTML +=
            " <li>\n" +
            "                <img src=\"\" alt=\"\">\n" +
            "                <div>\n" +
            "                    <h2>" + usersArray[i].name + "</h2>\n" +
            "                    <h3>\n" +
            "                        <span class=\"status green\"></span>\n" +
            "                        online\n" +
            "                    </h3>\n" +
            "                </div>\n" +
            "            </li>"

        sessionKey = usersArray[i].sessionKey.slice(usersArray[i].sessionKey.search("_") + 1)
        if (sessionId === sessionKey) {
            loggedInUserEmail = usersArray[i].email
            loggedInUserName = usersArray[i].name
        }
    }

}

function insertMessageHistory(msgList) {
    const msgDOMList = document.getElementById("chat");
    let msgTime;
    let youOrMeClass;
    let statusClass;
    for (let i = 0; i < msgList.length; i++) {
        msgTime = new Date(msgList[i].timestamp)
        youOrMeClass = msgList[i].email === loggedInUserEmail ? "me" : "you"
        statusClass = msgList[i].email === loggedInUserEmail ? "status blue" : "status green"

        msgDOMList.innerHTML +=
            "<li class=\"" + youOrMeClass + "\">\n" +
            "                <div class=\"entete\">\n" +
            "                    <span class=\"" + statusClass + "\"></span>\n" +
            "                    <h2>" + msgList[i].name + "</h2>\n" +
            "                    <h3>" + msgTime.toLocaleTimeString("en-us", {
                weekday: "long", year: "numeric", month: "short",
                day: "numeric", hour: "2-digit", minute: "2-digit"
            }) + "</h3>\n" +
            "                </div>\n" +
            "                <div class=\"triangle\"></div>\n" +
            "                <div class=\"message\">\n" +
            "                 " + msgList[i].data + "   \n" +
            "                </div>\n" +
            "            </li>"
    }
}

function insertNewMessage(msg) {
    const msgDOMList = document.getElementById("chat");
    let msgTime;
    let youOrMeClass;
    let statusClass;

    youOrMeClass = msg.email === loggedInUserEmail ? "me" : "you"
    statusClass = msg.email === loggedInUserEmail ? "status blue" : "status green"
    msgTime = new Date(msg.timestamp)

    msgDOMList.innerHTML +=
        "<li class=\"" + youOrMeClass + "\">\n" +
        "                <div class=\"entete\">\n" +
        "                    <span class=\"" + statusClass + "\"></span>\n" +
        "                    <h2>" + msg.name + "</h2>\n" +
        "                    <h3>" + msgTime.toLocaleTimeString("en-us", {
            weekday: "long", year: "numeric", month: "short",
            day: "numeric", hour: "2-digit", minute: "2-digit"
        }) + "</h3>\n" +
        "                </div>\n" +
        "                <div class=\"triangle\"></div>\n" +
        "                <div class=\"message\">\n" +
        "                 " + msg.data + "   \n" +
        "                </div>\n" +
        "            </li>"

    scrollChatScreenToBottom()

}


function insertRoomList(roomList) {
    const chatRoomList = document.getElementById("chat-room-select");
    for (let i = 0; i < roomList.length; i++) {
        const option = document.createElement("option");
        option.text = roomList[i].name;
        chatRoomList.add(option);
    }
}

function sendChatMsg() {
    let msgBox = document.getElementById("chat-msg-box")
    let msg = msgBox.value
    if (msg === "") return;
    webSocketCommand(WsCommands.MESSAGE, msg)
    msgBox.value = ""
}

function scrollChatScreenToBottom() {
    let chatScreen = document.getElementById("chat")
    //scroll chat screen to bottom showing new message
    chatScreen.scrollTop = chatScreen.scrollHeight;
}