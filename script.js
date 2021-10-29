let box = document.getElementById('firstspace')
textarea = document.getElementsByTagName('textarea')[0];
// Создаёт WebSocket - подключение.
const socket = new WebSocket('ws://localhost:8081');

// Соединение открыто
socket.addEventListener('open', function (event) {
    socket.send('Hello Server!');
});

// Наблюдает за сообщениями
socket.addEventListener('message', function (event) {
    console.log('Message from server ', event.data);
});
function ping(){
    textarea.value += box.value+"\n"
    
    socket.send("axaxaxa");
  }