// init
let conn, userId, roomId;
let form, msg, log, btnUser, btnRoom;

let lastData = {
  user: "",
  time: "",
  body: "",
};

const promptUserMsg = "Sign in as:";
const promptRoomMsg = "Join room: (blank to create a new one)";

window.onload = function () {
  let queryString = window.location.search;
  let params = new URLSearchParams(queryString);
  userId = params.get("user");
  roomId = params.get("room");

  form = document.getElementById("form");
  msg = document.getElementById("form-msg");
  log = document.getElementById("log");
  btnUser = document.getElementById("btn-user");
  btnRoom = document.getElementById("btn-room");

  setEvents();
  if (userId == null || !userId.trim()) {
    promptUserId();
  }
  if (roomId == null || !roomId.trim()) {
    promptRoomId();
  }
};

// trim string
let trim = (str) => {
  return str.replace(/ /g, "");
};

// https://cdn.jsdelivr.net/npm/nanoid/nanoid.js
let nanoid = (t = 21) => {
  let e = "",
    r = crypto.getRandomValues(new Uint8Array(t));
  for (; t--; ) {
    let n = 63 & r[t];
    e +=
      n < 36
        ? n.toString(36)
        : n < 62
        ? (n - 26).toString(36).toUpperCase()
        : n < 63
        ? "_"
        : "-";
  }
  return e;
};

// append log
let writeLog = (item) => {
  let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
  log.appendChild(item);
  if (doScroll) {
    log.scrollTop = log.scrollHeight - log.clientHeight;
  }
};

// clear log
let clearLog = () => {
  log.innerHTML = "";
};

// prompt user
let promptUserId = () => {
  userId = trim(prompt(promptUserMsg));
  if (userId == null || !userId) {
    promptUserId();
  } else {
    btnUser.innerHTML = `<b>${userId}</b>`;
  }
};

// prompt room
let promptRoomId = () => {
  newRoomId = !roomId
    ? trim(prompt(promptRoomMsg))
    : trim(prompt(promptRoomMsg, roomId));

  let isSameRoom = newRoomId == roomId;
  let isNewRoomValid = !!newRoomId && !isSameRoom;

  roomId = isNewRoomValid ? newRoomId : nanoid(8);
  if (!isSameRoom) {
    document.title = `#${roomId} - ChipBird IM`;
    btnRoom.innerHTML = `&#9998; <b>${roomId}</b>`;
    clearLog();
    openWs();
  }
};

// connect to websocket
let openWs = () => {
  if (window["WebSocket"]) {
    let wss = `wss://${document.location.host}/ws?room=${roomId}&user=${userId}`;
    conn = new WebSocket(wss);

    conn.onclose = function () {
      let item = document.createElement("div");
      item.classList.add("text-align-c");
      item.innerHTML = "<b>Connection closed.</b>";
      writeLog(item, log);
    };

    conn.onmessage = function (evt) {
      let messages = evt.data.split("\n");

      for (let i = 0; i < messages.length; i++) {
        let data = JSON.parse(messages[i]);
        let item = document.createElement("div");

        item.innerHTML =
          data.time == lastData.time && data.user == lastData.user
            ? ""
            : `<br><b>${data.user}</b> ` +
              `<span class='font-size-s'>${data.time}</span><br>`;
        item.innerHTML += `${data.body}`;

        lastData = data;
        writeLog(item);
      }
    };
  } else {
    let item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    writeLog(item);
  }
};

// set element event
let setEvents = () => {
  btnUser.onclick = function () {
    //promptUserId();
  };

  btnRoom.onclick = function () {
    promptRoomId();
  };

  form.onsubmit = function () {
    msg.value = msg.value.trim();
    if (conn && msg.value) {
      conn.send(msg.value);
      msg.value = "";
    }
    return false;
  };
};
