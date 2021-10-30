package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"

	"github.com/sparrc/go-ping"
)

func getping(conn net.Conn, host string, count int) {
	pinger, err := ping.NewPinger(host)
	pinger.Count = count
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
		time.Sleep(100 * time.Millisecond)
		out = fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		conn.Write([]byte(out))
		time.Sleep(100 * time.Millisecond)

		out = fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		conn.Write([]byte(out))
		time.Sleep(100 * time.Millisecond)

	}
	out := fmt.Sprintf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	conn.Write([]byte(out))
	time.Sleep(100 * time.Millisecond)

	pinger.Run()
	//conn.Write([]byte("\n"))

}
func gettrace(conn net.Conn, host string) {
	dateCmd := exec.Command("traceroute", host)
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	start := 0
	for i := 0; i < len(dateOut); i++ {
		if dateOut[i] == byte(10) {
			//fmt.Println(string(dateOut[start:i]))
			conn.Write([]byte(string(dateOut[start:i]) + "\n"))
			start = i + 1
			time.Sleep(100 * time.Millisecond)
		}
	}

	//fmt.Println(string(dateOut))
}
func decode(str string) (mode, host string, count int) {
	var num []int
	num = append(num, len(str)-1)
	num = append(num, len(str)-1)
	j := 0
	for i := 0; i < len(str); i++ {
		if str[i] == byte(44) {
			num[j] = i
			j++
		}

	}
	count = 0
	if num[1] != len(str)-1 {
		count, _ = strconv.Atoi(str[(num[1] + 1) : len(str)-1])
	}

	return str[0:num[0]], str[(num[0] + 1):num[1]], count
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
		fmt.Print(string(message))
		// Процесс выборки для полученной строки
		// Отправить новую строку обратно клиенту

		mode, host, count := decode(message)
		switch mode {

		case "ping":
			go getping(conn, host, count)
		case "trace":
			go gettrace(conn, host)
		}

	}
}
func main() {
	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", ":8081")
	for {
		conn, err := ln.Accept()
		if err == nil {
			go client(conn)
		}
	}
}
