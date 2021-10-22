package main

import (
	"fmt"

	"github.com/skorobogatov/input"
	"github.com/sparrc/go-ping"
)

func main() {
	fmt.Println("Введите хост")
	var host string
	input.Scanf("%s", &host)
	pinger, err := ping.NewPinger(host)
	fmt.Print("Введите кол-во запросов, либо оставьте 0")
	var count int
	input.Scanf("%d", &count)
	if count > 0 {
		pinger.Count = count
	}

	if err != nil {

		fmt.Printf("ERROR: %s\n", err.Error())

		return

	}

	pinger.OnRecv = func(pkt *ping.Packet) {

		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",

			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		//fmt.Println(pkt)

	}

	pinger.OnFinish = func(stats *ping.Statistics) {

		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)

		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",

			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)

		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",

			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	pinger.Run()
}
