const app = require('express')();
const http = require('http').Server(app);
const io = require('socket.io')(http);
const port = process.env.PORT || 3000;

app.get('/', (req, res) => {
  res.sendFile(__dirname + '/views/index.html');
});

app.get('/room', (req, res) => {
  res.sendFile(__dirname + '/views/room.html');
});

io.on('connection', (socket) => {

  let username;
  let room;

  socket.on('new user', function(name, newroom){
    // join new room and tell everyone you've entered
    room = newroom;
    socket.join(newroom);
    username = name;
    socket.broadcast.to(room).emit('chat sub msg', `${name} has joined the chat`);
  });

  // sends message to everyone in the chat
  socket.on('chat message', (msg) => {
    io.to(room).emit('chat message', `${username}: ${msg}`);
  });

  // events to trigger when socket disconnects
  socket.on('disconnect', () => {
    // let everyone know you've left the chat
    socket.broadcast.to(room).emit('chat sub msg', `${username} has left the chat`);
  });
  
});


http.listen(port, () => {
  console.log(`Socket.IO server running at http://localhost:${port}/`);
});
