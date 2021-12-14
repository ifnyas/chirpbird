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

const wsConnectingMsg = "Connecting...";
const wsOpenedMsg = "You are connected!";
const wsClosedMsg = "Connection closed.";
const wsErrMsg = "Your browser does not support WebSockets.";

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
let appendLog = (item) => {
  let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
  log.appendChild(item);
  if (doScroll) {
    log.scrollTop = log.scrollHeight - log.clientHeight;
  }
};

// clear log
let clearLog = () => {
  log.innerHTML = "";
  lastData = {
    user: "",
    time: "",
    body: "",
  };
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
    eventWsView(wsConnectingMsg);
    startWs();
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

// connect to websocket
let startWs = () => {
  if (window["WebSocket"]) {
    let protocol = window.location.protocol == "https:" ? "wss" : "ws";
    let wss = `${protocol}://${window.location.host}/ws?key=${roomId}%3A${userId}`;

    conn = new WebSocket(wss);
    conn.onopen = () => {
      openWsView();
    };
    conn.onclose = () => {
      eventWsView(wsClosedMsg);
    };
    conn.onmessage = (ev) => {
      rcvWsView(JSON.parse(ev.data));
    };
  } else {
    eventWsView(wsErrMsg);
  }
};

let openWsView = async () => {
  await loadLogs();
  msg.disabled = false;
};

let eventWsView = (ev) => {
  let item = document.createElement("div");
  item.classList.add("text-align-c");
  item.innerHTML = `<br>${ev}<br>`;
  appendLog(item, log);
};

let rcvWsView = (data) => {
  let item = document.createElement("div");
  item.innerHTML =
    data.time == lastData.time && data.user == lastData.user
      ? ""
      : `<br><b>${data.user}</b> ` +
        `<span class='font-size-s'>${data.time}</span><br>`;
  item.innerHTML += `${data.body}`;

  lastData = data;
  appendLog(item);
};

async function loadLogs() {
  let url = `${window.location.protocol}/logs?room=${roomId}`;
  let obj = await (await fetch(url)).json();

  obj.data.forEach((item) => {
    rcvWsView(item);
  });
}
