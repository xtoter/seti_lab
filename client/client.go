package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	ln, _ := net.Listen("tcp", ":8081")
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		message, err1 := bufio.NewReader(conn).ReadString('\n')
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		conn.Write([]byte(string("axaxax") + "\n"))
		// Распечатываем полученое сообщение
		fmt.Print(string(message))
	}

}
