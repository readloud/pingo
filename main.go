package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/hideckies/pingo/sub"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type Pingo struct {
	Count    int
	Host     string
	Interval time.Duration
	Packet   *sub.Packet

	// Channel
	done chan interface{}
	mtx  sync.Mutex
}

// Execute ping
// Reference: https://pkg.go.dev/golang.org/x/net@v0.0.0-20221004154528-8021a29435af/icmp#example-PacketConn-NonPrivilegedPing
func (p *Pingo) Run(statistics *sub.Statistics) error {
	// packetconn, err := icmp.ListenPacket(p.Packet.Network, p.Packet.SrcAddr.String())
	packetconn, err := sub.ListenPacket(p.Packet.Network, p.Packet.SrcAddr.String())
	if err != nil {
		log.Fatalf("ICMP ListenPacket Error: %v\n", err)
	}
	defer packetconn.Close()

	pktconn := packetconn.IPv4PacketConn()
	if err := pktconn.SetTTL(p.Packet.TTL); err != nil {
		log.Fatalf("SetTTL Error: %v\n", err)
	}

	c := 1
	for range time.Tick(p.Interval) {
		p.Packet.Seq = c
		body := &icmp.Echo{
			ID:   p.Packet.ID,
			Seq:  p.Packet.Seq,
			Data: []byte(p.Packet.Data),
		}
		msg := &icmp.Message{
			Type: p.Packet.ICMPType,
			Code: 0,
			Body: body,
		}

		wb, err := msg.Marshal(nil)
		if err != nil {
			log.Fatalf("Marshal Error: %v\n", err)
		}

		var dst net.Addr
		if strings.Contains(p.Packet.Network, "udp") {
			dst = &net.UDPAddr{IP: p.Packet.DestAddr.IP, Zone: p.Packet.DestAddr.Zone}
		} else {
			dst = p.Packet.DestAddr
		}
		if _, err := packetconn.WriteTo(wb, dst); err != nil {
			log.Fatalf("WriteTo Error: %v\n", err)
		}

		rb := make([]byte, 1500)
		n, peer, err := packetconn.ReadFrom(rb)
		if err != nil {
			log.Fatalf("ReadFrom Error: %v\n", err)
		}

		rm, err := icmp.ParseMessage(p.Packet.ProtoNum, rb[:n])
		if err != nil {
			log.Fatalf("ParseMessage Erorr: %v", err)
		}

		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			statistics.Result(peer, n, *p.Packet)
			statistics.Received++
		case ipv6.ICMPTypeEchoReply:
			statistics.Result(peer, n, *p.Packet)
			statistics.Received++
		default:
			fmt.Printf(":-< faled %+v\n", rm)
			statistics.Loss++
		}

		c++
		statistics.Transmitted++

		if p.Count != 0 && c > p.Count {
			break
		}
	}

	return nil
}

func (p *Pingo) Stop(statistics *sub.Statistics) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	open := true
	select {
	case _, open = <-p.done:
	default:
	}

	if open {
		close(p.done)
	}

	// display the result
	statistics.FinalResult()

	os.Exit(0)
}

func NewPingo(flag *sub.Flag, packet *sub.Packet) *Pingo {
	var p Pingo
	p.Count = flag.Count
	p.Host = flag.Target
	p.Packet = packet
	p.done = make(chan interface{})

	interval, err := time.ParseDuration(flag.Interval + "s")
	if err == nil {
		p.Interval = interval
	} else {
		fmt.Println(sub.ErrIncorrectValueInterval)
		p.Interval = 1 * time.Second
	}

	return &p
}

func main() {
	var f sub.Flag

	err := f.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	packet := sub.NewPacket(&f)
	pingo := NewPingo(&f, packet)
	statistics := sub.NewStatistics(&f, packet)

	statistics.Start(pingo.Host, pingo.Packet.DestAddr.String())

	// Listen for Ctrl+c signal
	cch := make(chan os.Signal, 1)
	signal.Notify(cch, os.Interrupt)
	go func() {
		for range cch {
			pingo.Stop(statistics)
		}
	}()

	err = pingo.Run(statistics)
	if err != nil {
		fmt.Println("Error pingo")
	}

	statistics.FinalResult()
}
