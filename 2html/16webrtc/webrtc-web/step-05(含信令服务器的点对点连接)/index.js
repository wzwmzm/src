'use strict';
//本程序为信令服务器
var os = require('os');
var nodeStatic = require('node-static');
var http = require('http');
var socketIO = require('socket.io');

//网页服务, 端口8080
var fileServer = new(nodeStatic.Server)();
var app = http.createServer(function(req, res) {
  fileServer.serve(req, res);
}).listen(8080);
console.log("http://localhost:8080");

//websocket
var io = socketIO.listen(app);
io.sockets.on('connection', function(socket) {

  // convenience function to log server messages on the client
  function log() {
    var array = ['Message from server:'];
    array.push.apply(array, arguments);
      //.push()是给数组array增加元素arguments
      //.apply()是指定作用域为array
      //arguments是函数的隐式参数
    socket.emit('log', array);
    
    console.log(arguments);//<---------服务器端的提示
  }

    //简单地广播, 需要扩展成对特定房间的转发
  socket.on('message', function(message) {
    log('Client said: ', message);
    // for a real app, would be room-only (not broadcast)
    socket.broadcast.emit('message', message);
  });

    //
  socket.on('create or join', function(room) {//room为客户端手动输入
    log('Received request to create or join room ' + room);

    var clientsInRoom = io.sockets.adapter.rooms[room];
      //room如果存在则返回rooms; 如果不存在则返回false
    var numClients = clientsInRoom ? Object.keys(clientsInRoom.sockets).length : 0;
    log('Room ' + room + ' now has ' + numClients + ' client(s)');

    if (numClients === 0) {//房间没人,则建立房间
      socket.join(room);//加入房间
      log('Client ID ' + socket.id + ' created room ' + room);
      socket.emit('created', room, socket.id);

    } else if (numClients === 1) {//房间只有一个人,则加入
      log('Client ID ' + socket.id + ' joined room ' + room);
      io.sockets.in(room).emit('join', room);
      socket.join(room);//加入房间
      socket.emit('joined', room, socket.id);
      io.sockets.in(room).emit('ready');
    } else { // max two clients, 有两个人就满了
      socket.emit('full', room);
    }
  });

    //发送服务器的IPv4地址(非127.0.0.1)给客户端
  socket.on('ipaddr', function() {
    var ifaces = os.networkInterfaces();
    for (var dev in ifaces) {
      ifaces[dev].forEach(function(details) {
        if (details.family === 'IPv4' && details.address !== '127.0.0.1') {
          socket.emit('ipaddr', details.address);
        }
      });
    }
  });

  socket.on('bye', function(){
    console.log('received bye');
  });

});
