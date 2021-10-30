package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/", websocket.Handler(handler))
	http.ListenAndServe("localhost:3000", nil)
}

func handler(c *websocket.Conn) {
	var s string
	fmt.Fscan(c, &s)
	fmt.Println("Received:", s)
	fmt.Fprint(c, "How do you do?")
}

// test using these in browser console:
// var sock = new WebSocket("ws://localhost:3000/");
// sock.onmessage = function(m) { console.log("Received:", m.data); }
// sock.send("Hello!\n")
