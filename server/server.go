package main

import (
	"fmt"
	"os/exec"

	"github.com/sparrc/go-ping"
)

func main() {
	mode := "trace"
	switch mode {
	case "ping":
		pinger, err := ping.NewPinger("ya.ru")
		//pinger.Count = count
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
	case "trace":
		dateCmd := exec.Command("traceroute", "ya.ru")
		dateOut, err := dateCmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dateOut))
	}

}
