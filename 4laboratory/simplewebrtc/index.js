


var webrtc = new SimpleWebRTC({

    // the id/element dom element that will hold "our" video

    localVideoEl: 'localVideo',

    // the id/element dom element that will hold remote videos

    remoteVideosEl: 'remoteVideos',

    // immediately ask for camera access

    autoRequestMedia: true,

    url:'http://qq.openauth.me:8888',  //signal服务器

    nick: 'yubaolee'

});







// we have to wait until it's ready

webrtc.on('readyToCall', function () {

    // you can name it anything

    webrtc.joinRoom('room1');



    // Send a chat message

    $('#send').click(function () {

        var msg = $('#text').val();

        webrtc.sendToAll('chat', { message: msg, nick: webrtc.config.nick });

        $('#messages').append('<br>You:<br>' + msg + '\n');

        $('#text').val('');

    });

});



//For Text Chat ------------------------------------------------------------------

// Await messages from others

webrtc.connection.on('message', function (data) {

    if (data.type === 'chat') {

        console.log('chat received', data);

        $('#messages').append('<br>' + data.payload.nick + ':<br>' + data.payload.message + '\n');

    }

});

