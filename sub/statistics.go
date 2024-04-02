package sub

import (
	"fmt"
	"net"
	"strings"

	"github.com/fatih/color"
)

type Statistics struct {
	Proto string

	Transmitted int
	Received    int
	Loss        int

	NoColor bool
}

func (s *Statistics) Start(host, dstaddr string) {
	fmt.Printf("%s Ping %s (%s)\n", strings.ToUpper(s.Proto), host, dstaddr)
	fmt.Println()
}

func (s *Statistics) Result(peer net.Addr, n int, packet Packet) {
	if s.NoColor {
		fmt.Printf(":-) from=%v::bytes=%d::id=0x%x::seq=%d::ttl=%d::protocol=%s\n", peer, n, packet.ID, packet.Seq, packet.TTL, s.Proto)
	} else {
		color.HiCyan(":-) from=%v::bytes=%d::id=0x%x::seq=%d::ttl=%d::protocol=%s\n", peer, n, packet.ID, packet.Seq, packet.TTL, s.Proto)
	}
}

func (s *Statistics) FinalResult() {
	fmt.Println()
	fmt.Println("--- statistics ---------------------------------------")

	if s.NoColor {
		fmt.Printf("transmitted=%d::received=%d::loss=%d::protocol=%s\n", s.Transmitted, s.Received, s.Loss, s.Proto)
	} else {
		transmitted := color.GreenString("transmitted=%d", s.Transmitted)
		received := color.CyanString("received=%d", s.Received)
		loss := color.RedString("loss=%d", s.Loss)
		proto := color.YellowString("protocol=%s", s.Proto)
		fmt.Printf("%s::%s::%s::%s\n", transmitted, received, loss, proto)
	}

	fmt.Println("------------------------------------------------------")
}

func NewStatistics(f *Flag, p *Packet) *Statistics {
	proto := "icmp"
	if strings.Contains(p.Network, "udp") {
		proto = "udp"
	}

	return &Statistics{
		Proto:       proto,
		Transmitted: 0,
		Received:    0,
		Loss:        0,
		NoColor:     f.NoColor,
	}
}
