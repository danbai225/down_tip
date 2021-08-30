var socket = null;
var pingFlag =false
let tagMsgMap=new Map();
//初始化websocket
function initWebSocket() {
  //ws地址
  var wsUri = "ws://127.0.0.1:7989/socket/ws";
  socket = new WebSocket(wsUri);
  socket.onmessage = function(e) {
    onMessage(e);
  };
  socket.onclose = function(e) {
    onClose(e);
  };
  socket.onopen = function(e) {
    onOpen(e);
  };

  //连接发生错误的回调方法
  socket.onerror = function() {
    console.log("WebSocket连接发生错误");
  };
  if (!pingFlag){
    setTimeout(function(){
      sendMsg({"type":"ping","data":Date.parse(new Date())})
    }, 1000);
    pingFlag=true
  }
  
}
// 实际调用的方法
function sendMsg(data) {

  if (socket.readyState === socket.OPEN) {
    //若是ws开启状态
    onSend(data);
  } else if (socket.readyState === socket.CONNECTING) {
    // 若是 正在开启状态，则等待1s后重新调用
    setTimeout(function() {
      sendMsg(data);
    }, 1000);
  } else {
    initWebSocket()
    // 若未开启 ，则等待1s后重新调用
    setTimeout(function() {
      sendMsg(data);
    }, 1000);
  }
}
function setTagMsg(tag,func){
  tagMsgMap.set(tag,func);
}

//数据接收
function onMessage(e) {
  console.log(e.data);
  var data=JSON.parse(e.data);
  if (tagMsgMap.has(data.tag)){
    tagMsgMap.get(data.tag)(data)
  }
}

//数据发送
function onSend(agentData) {
  socket.send(JSON.stringify(agentData));
}

//关闭
function onClose(e) {
  console.log("connection closed (" + e.code + ")");
}

function onOpen(e) {
  console.log("连接成功",e);
}

export { sendMsg ,initWebSocket,setTagMsg};
