'use strict';
//本程序在local和remote都会各自执行, 所以它们有公共部分和私有部分
var isChannelReady = false;//两个人都加入同一个房间, isChannelReady = true
var isInitiator = false;//只有发起人建立房间并进入房间后, isInitiator = true;
//或者说发起人 isInitiator = true   , 加入者 isInitiator = false;
var isStarted = false;//PeerConnection建立后 isStarted = true  
var localStream;
var pc;//虽然本程序只有一个PeerConnection,但由于在local和remote都会各自执行,所以其实有两个pc
var remoteStream;
var turnReady;//turn中继服务器准备好

//stun是打洞服务器,  'stun:stun.l.google.com:19302'
//turn是中继服务器,  'https://computeengineondemand.appspot.com/turn?username=41784574&key=4080218913'    ;超出5g要付费
var pcConfig = {
  'iceServers': [{
    'urls': 'stun:stun.l.google.com:19302'
  }]
};

// Set up audio and video regardless of what devices are present.
var sdpConstraints = {
  offerToReceiveAudio: true,
  offerToReceiveVideo: true
};

////////websockets 实现房间管理功能////////////////////////////

var room = 'foo';
// Could prompt for room name:
// room = prompt('Enter room name:');

var socket = io.connect();

if (room !== '') {
  socket.emit('create or join', room);
	//两个客户端都会发送此消息,但只有快的客户端建立了房间,慢的则加入房间
  console.log('Attempted to create or  join room', room);
}

socket.on('created', function(room) {//只有快的客户端会收到此消息
  console.log('Created room ' + room);
  isInitiator = true;//<----只有快的客户端会
});

socket.on('full', function(room) {
  console.log('Room ' + room + ' is full');
});

socket.on('join', function (room){//只有慢的客户端会收到此消息
  console.log('Another peer made a request to join room ' + room);
  console.log('This peer is the initiator of room ' + room + '!');
  isChannelReady = true;
});

socket.on('joined', function(room) {//只有慢的客户端会收到此消息
  console.log('joined: ' + room);
  isChannelReady = true;
});

socket.on('log', function(array) {
  console.log.apply(console, array);
});

/////////实现信令服务器功能///////////////////////////////////////

function sendMessage(message) {
  console.log('Client sending message: ', message);
  socket.emit('message', message);
}

// This client receives a message
socket.on('message', function(message) {
  console.log('Client received message:', message);
  if (message === 'got user media') {//两个客户端都会收到
    maybeStart();//建立PeerConnection, 包括绑定事件处理函数, 主要是发送icecandidate
  } else if (message.type === 'offer') {                     //remotePeer
	  //发起人还没建立并进入房间(或者说加入者) && PeerConnection没有建立
	  if (!isInitiator && !isStarted) {
      maybeStart();
    }
    pc.setRemoteDescription(new RTCSessionDescription(message));
    doAnswer();
  } else if (message.type === 'answer' && isStarted) {      //localPeer
    pc.setRemoteDescription(new RTCSessionDescription(message));
  } else if (message.type === 'candidate' && isStarted) {   //icecandidate
    var candidate = new RTCIceCandidate({
      sdpMLineIndex: message.label,
      candidate: message.candidate
    });
    pc.addIceCandidate(candidate);
  } else if (message === 'bye' && isStarted) {              //结束
    handleRemoteHangup();
  }
});

///////////主  程  序/////////////////////////////////////////

var localVideo = document.querySelector('#localVideo');
var remoteVideo = document.querySelector('#remoteVideo');

navigator.mediaDevices.getUserMedia({
  audio: false,
  video: true
})
.then(gotStream)
.catch(function(e) {
  alert('getUserMedia() error: ' + e.name);
});

function gotStream(stream) {
  console.log('Adding local stream.');
  localStream = stream;
  localVideo.srcObject = stream;
  sendMessage('got user media');
  if (isInitiator) {
    maybeStart();
  }
}

var constraints = {
  video: true
};

console.log('Getting user media with constraints', constraints);

if (location.hostname !== 'localhost') {
  requestTurn(//请求中继服务器,在这里可以设置中继服务器 (注意:因为获取不到,可能会出错)
    'https://computeengineondemand.appspot.com/turn?username=41784574&key=4080218913'
  );
}

function maybeStart() {//建立PeerConnection, 包括绑定事件处理函数, 主要是发送icecandidate
  console.log('>>>>>>> maybeStart() ', isStarted, localStream, isChannelReady);
  if (!isStarted && typeof localStream !== 'undefined' && isChannelReady) {
    console.log('>>>>>> creating peer connection');
    createPeerConnection();//建立PeerConnection, 包括绑定事件处理函数, 主要是发送icecandidate
    pc.addStream(localStream);//添加媒体流
    isStarted = true;
    console.log('isInitiator', isInitiator);
    if (isInitiator) {
      doCall();
    }
  }
}

window.onbeforeunload = function() {
  sendMessage('bye');
};

/////////////////////////////////////////////////////////
////PeerConnection, 在local和remot都会建立(执行). 建立后 isStarted = true
function createPeerConnection() {//建立PeerConnection, 包括绑定事件处理函数
  try {
    pc = new RTCPeerConnection(null);//<-------
    pc.onicecandidate = handleIceCandidate;//发送icecandidate
    pc.onaddstream = handleRemoteStreamAdded;//添加媒体流
    pc.onremovestream = handleRemoteStreamRemoved;//纯log 
    console.log('Created RTCPeerConnnection');
  } catch (e) {
    console.log('Failed to create PeerConnection, exception: ' + e.message);
    alert('Cannot create RTCPeerConnection object.');
    return;
  }
}

function handleIceCandidate(event) {//发送icecandidate
  console.log('icecandidate event: ', event);
  if (event.candidate) {
    sendMessage({
      type: 'candidate',
      label: event.candidate.sdpMLineIndex,
      id: event.candidate.sdpMid,
      candidate: event.candidate.candidate
    });
  } else {
    console.log('End of candidates.');
  }
}

function handleCreateOfferError(event) {
  console.log('createOffer() error: ', event);
}

//doCall 和 doAnswer, createOffer 和 createAnswer 本质和内容基本一样
//1,设置本地SDP; 2,开始收集icecandidate; 3,发送SDP给对方
function doCall() {
  console.log('Sending offer to peer');
  pc.createOffer(setLocalAndSendMessage, handleCreateOfferError);
}

function doAnswer() {
  console.log('Sending answer to peer.');
  pc.createAnswer().then(
    setLocalAndSendMessage,
    onCreateSessionDescriptionError
  );
}

function setLocalAndSendMessage(sessionDescription) {
  pc.setLocalDescription(sessionDescription);//<--------开始收集icecandidate
  console.log('setLocalAndSendMessage sending message', sessionDescription);
  sendMessage(sessionDescription);//<-------------------发送SDP
}

function onCreateSessionDescriptionError(error) {
  trace('Failed to create session description: ' + error.toString());
	//console.trace() 
}

//取得 TURN 中继服务器地址, 可以在两个地方设置:(注意:因为获取不到,可能会出错)
//1, 在 pcConfig 中设置
//2, requestTurn(turnURL) 中的 turnURL 设置
function requestTurn(turnURL) {//(注意:因为获取不到,可能会出错)
  var turnExists = false;
  for (var i in pcConfig.iceServers) {//从pcConfig中提取
    if (pcConfig.iceServers[i].urls.substr(0, 5) === 'turn:') {
      turnExists = true;
      turnReady = true;
      break;
    }
  }
  if (!turnExists) {
    console.log('Getting TURN server from ', turnURL);
    // No TURN server. Get one from computeengineondemand.appspot.com:
    var xhr = new XMLHttpRequest();//在turnURL中设置,最终会加入到pcConfig中
    xhr.onreadystatechange = function() {
      if (xhr.readyState === 4 && xhr.status === 200) {
        var turnServer = JSON.parse(xhr.responseText);
        console.log('Got TURN server: ', turnServer);
        pcConfig.iceServers.push({
          'urls': 'turn:' + turnServer.username + '@' + turnServer.turn,
          'credential': turnServer.password
        });
        turnReady = true;
      }
    };
    xhr.open('GET', turnURL, true);
    xhr.send();
  }
}

function handleRemoteStreamAdded(event) {
  console.log('Remote stream added.');
  remoteStream = event.stream;
  remoteVideo.srcObject = remoteStream;
}

function handleRemoteStreamRemoved(event) {
  console.log('Remote stream removed. Event: ', event);
}

function hangup() {
  console.log('Hanging up.');
  stop();
  sendMessage('bye');
}

function handleRemoteHangup() {
  console.log('Session terminated.');
  stop();
  isInitiator = false;
}

function stop() {
  isStarted = false;
  pc.close();
  pc = null;
}
