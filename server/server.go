package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"

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
		out := fmt.Sprintf("\n--- %s ping statistics ---\n", stats.Addr)
		conn.Write([]byte(out))
		out = fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		conn.Write([]byte(out))

		out = fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		conn.Write([]byte(out))

	}
	out := fmt.Sprintf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	conn.Write([]byte(out))

	pinger.Run()

}
func gettrace(conn net.Conn) {
	dateCmd := exec.Command("traceroute", "ya.ru")
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	conn.Write([]byte(string(dateOut) + "/n"))
	fmt.Println(string(dateOut))
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

		mode := message
		fmt.Print([]byte(message))
		switch mode {

		case "ping":
			getping(conn)
		case "trace" + string(byte(10)):
			fmt.Println("d")
			gettrace(conn)
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
