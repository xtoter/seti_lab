let box = document.getElementById('firstspace')
textarea = document.getElementsByTagName('textarea')[0];
var socket = new WebSocket("ws://localhost:3000/");
socket.onmessage = function(m) { textarea.value += m.data+"\n"; }

function trace() {
  console.log("trace,"+box.value)
  socket.send("trace,"+box.value);
}
function ping(){
  console.log("ping,"+box.value)
    socket.send("ping,"+box.value);
}