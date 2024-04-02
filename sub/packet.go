package sub

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const (
	PROTO_NUM_ICMP_IPv4 = 1
	PROTO_NUM_ICMP_IPv6 = 58
)

type Packet struct {
	Data     string
	ICMPType icmp.Type
	ID       int
	Network  string
	Proto    string
	ProtoNum int
	Seq      int
	TTL      int

	SrcAddr  *net.IPAddr
	DestAddr *net.IPAddr
}

// Create a new Packet
func NewPacket(f *Flag) *Packet {
	proto := getProtocol(f.Unprivileged, f.UseIPv6)
	protoNum := getProtocolNumber(f.UseIPv6)
	srcIp := getSrcIP(f.UseIPv6)

	return &Packet{
		Data:     f.Data,
		ICMPType: getICMPType(f.UseIPv6),
		ID:       getID(),
		Network:  getNetwork(proto, protoNum),
		Proto:    proto,
		ProtoNum: protoNum,
		Seq:      0,
		TTL:      f.TTL,
		SrcAddr:  resolve(proto, srcIp),
		DestAddr: resolve(proto, f.Target),
	}
}

// Get the ICMP type
func getICMPType(useIpv6 bool) icmp.Type {
	var icmpType icmp.Type

	if useIpv6 {
		icmpType = ipv6.ICMPTypeEchoRequest
	} else {
		icmpType = ipv4.ICMPTypeEcho
	}

	return icmpType
}

// Get the packet ID
func getID() int {
	var seed int64 = time.Now().UnixNano()
	newseed := atomic.AddInt64(&seed, 1)

	r := rand.New(rand.NewSource(newseed))
	return r.Intn(math.MaxUint16)
}

// Get the network
func getNetwork(proto string, protoNum int) string {
	var network string

	if strings.Contains(proto, "udp") {
		network = proto
	} else {
		network = proto + ":" + strconv.Itoa(protoNum)
	}

	return network
}

// Get the protocol name
func getProtocol(unprivileged, useIPv6 bool) string {
	var proto string

	if unprivileged && useIPv6 {
		proto = "udp6"
	} else if unprivileged && !useIPv6 {
		proto = "udp4"
	} else if !unprivileged && useIPv6 {
		proto = "ip6"
	} else if !unprivileged && !useIPv6 {
		proto = "ip4"
	}

	return proto
}

// Get the protocol number
func getProtocolNumber(useIPv6 bool) int {
	var protoNum int

	if useIPv6 {
		protoNum = PROTO_NUM_ICMP_IPv6
	} else {
		protoNum = PROTO_NUM_ICMP_IPv4
	}

	return protoNum
}

// Get the source IP
func getSrcIP(useIPv6 bool) string {
	var srcIp string

	if useIPv6 {
		srcIp = "::"
	} else {
		srcIp = "0.0.0.0"
	}

	return srcIp
}

// Resolve IP address
func resolve(proto, address string) *net.IPAddr {
	var network string

	if strings.Contains(proto, "udp") {
		network = strings.Replace(proto, "udp", "ip", 1)
	} else {
		network = proto
	}

	addr, err := net.ResolveIPAddr(network, address)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return addr
}
