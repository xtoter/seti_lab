package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func read(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err == nil {
			fmt.Print(message)
		}
	}

}
func main2() {

	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	go read(conn)
	for {
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// Отправляем в socket
		fmt.Fprintf(conn, text+"\n")
		// Прослушиваем ответ
	}
}
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
