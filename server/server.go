package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/sparrc/go-ping"
)

func getping(conn net.Conn) {
	pinger, err := ping.NewPinger("ya.ru")
	//pinger.Count = count
	if err != nil {
		conn.Write([]byte("Error: " + err.Error()))
		//fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	pinger.OnRecv = func(pkt *ping.Packet) {

		out := fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v\n",

			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		conn.Write([]byte(out))
		//fmt.Println(pkt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		out := fmt.Sprintf("--- %s ping statistics ---", stats.Addr)
		conn.Write([]byte(out))
		out = fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		conn.Write([]byte(out))

		out = fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		conn.Write([]byte(out))

	}
	out := fmt.Sprintf("PING %s (%s):", pinger.Addr(), pinger.IPAddr())
	conn.Write([]byte(out))

	pinger.Run()
	//conn.Write([]byte("\n"))

}
func gettrace(conn net.Conn) {
	dateCmd := exec.Command("traceroute", "ya.ru")
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	start := 0
	for i := 0; i < len(dateOut); i++ {
		if dateOut[i] == byte(10) {
			fmt.Println(string(dateOut[start:i]))
			conn.Write([]byte(string(dateOut[start:i]) + "\n"))
			start = i + 1
			time.Sleep(100 * time.Millisecond)
		}
	}

	//fmt.Println(string(dateOut))
}
func decode(str string) (mode, host string) {
	num := len(str) - 1
	for i := 0; i < len(str); i++ {
		if str[i] == byte(44) {
			num = i
		}

	}
	return str[0:num], str[(num + 1):]
}
func client(conn net.Conn) {
	for {
		// Будем прослушивать все сообщения разделенные \n
		message, err1 := bufio.NewReader(conn).ReadString('\n')
		if err1 != nil {
			fmt.Println("error")
			return
		}
		// Распечатываем полученое сообщение
		fmt.Print("Message Received:", string(message))
		// Процесс выборки для полученной строки
		// Отправить новую строку обратно клиенту

		mode, host := decode(message)
		host = host
		switch mode {

		case "ping":
			go getping(conn)
		case "trace":
			go gettrace(conn)
		}

	}
}
func main() {
	fmt.Println("Launching server...")
	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":8081")
	// Открываем порт
	for {
		conn, err := ln.Accept()
		if err == nil {
			go client(conn)
		}
	}
}
