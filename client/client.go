//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
var addr = flag.String("addr", "localhost:8080", "http service address")
var addrweb = flag.String("addrweb", "localhost:9090", "http service address")
var server *websocket.Conn

func read(c, d *websocket.Conn) {
	done := make(chan struct{})
	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		err = d.WriteMessage(websocket.TextMessage, message)
		//fmt.Printf("recv: %s\n", message)
	}
}
func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	fmt.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
	}
	defer c.Close()
	//go read(c)
	server = c
	//

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.ListenAndServe(*addrweb, nil)
	for {

	}
}
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		go read(server, c)
		server.WriteMessage(websocket.TextMessage, []byte(message))

		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
	var number = document.getElementById("number");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print(evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("ping").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("ping "+input.value+" "+number.value);
        ws.send("ping,"+input.value+","+number.value);
        return false;
    };
	document.getElementById("trace").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("trace "+input.value);
        ws.send("trace,"+input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="ya.ru">
<input id="number" type="number" name="t" value="2" min="0" max="9999" step="1">
<button id="ping">Ping</button>
<button id="trace">Trace</button>
</form>

<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td><td valign="top" width="50%">
</td></tr></table>
</body>
</html>
`))
