package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sparrc/go-ping"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func getping(c *websocket.Conn, host string, count int) {
	pinger, err := ping.NewPinger(host)
	pinger.Count = count
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
		//fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	pinger.OnRecv = func(pkt *ping.Packet) {

		out := fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v\n",

			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		c.WriteMessage(websocket.TextMessage, []byte(out))
		//fmt.Println(pkt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		out := fmt.Sprintf("\n--- %s ping statistics ---\n", stats.Addr)
		c.WriteMessage(websocket.TextMessage, []byte(out))
		time.Sleep(100 * time.Millisecond)
		out = fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		c.WriteMessage(websocket.TextMessage, []byte(out))
		time.Sleep(100 * time.Millisecond)

		out = fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		c.WriteMessage(websocket.TextMessage, []byte(out))
		time.Sleep(100 * time.Millisecond)

	}
	out := fmt.Sprintf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	c.WriteMessage(websocket.TextMessage, []byte(out))
	time.Sleep(100 * time.Millisecond)

	pinger.Run()
	//conn.Write([]byte("\n"))

}
func gettrace(c *websocket.Conn, host string) {
	dateCmd := exec.Command("traceroute", host)
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	start := 0
	for i := 0; i < len(dateOut); i++ {
		if dateOut[i] == byte(10) {
			//fmt.Println(string(dateOut[start:i]))
			c.WriteMessage(websocket.TextMessage, []byte(string(dateOut[start:i])+"\n"))
			start = i + 1
			time.Sleep(100 * time.Millisecond)
		}
	}

	//fmt.Println(string(dateOut))
}
func decode(str string) (mode, host string, count int) {
	var num []int
	num = append(num, len(str)-1)
	num = append(num, -1)
	j := 0
	for i := 0; i < len(str); i++ {
		if str[i] == byte(44) {
			num[j] = i
			j++
		}

	}
	count = 0
	if num[1] != -1 {
		count, _ = strconv.Atoi(str[(num[1] + 1):])
	} else {
		num[1] = len(str)
	}

	return str[0:num[0]], str[(num[0] + 1):num[1]], count
}
func client(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fmt.Println(string(message))

		// Распечатываем полученое сообщение
		// Процесс выборки для полученной строки
		// Отправить новую строку обратно клиенту

		mode, host, count := decode(string(message))
		switch mode {

		case "ping":
			go getping(c, host, count)
		case "trace":
			go gettrace(c, host)
		}

	}
}

var upgrader = websocket.Upgrader{} // use default options
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	client(c)
}
func main() {
	fmt.Println("Launching server...")
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
