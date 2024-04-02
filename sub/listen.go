package sub

import (
	"net"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const sysIP_STRIPHDR = 0x17 // for now only darwin supports this option

// A packet listener
// Reference: https://cs.opensource.google/go/x/net/+/8021a294:icmp/listen_posix.go;bpv=1;bpt=1;l=47?gsn=ListenPacket&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fnet%3Flang%3Dgo%3Fpath%3Dicmp%23func%2520ListenPacket
func ListenPacket(network, address string) (*PacketConn, error) {
	var family, proto int
	switch network {
	case "udp4":
		family, proto = syscall.AF_INET, IANA_ProtocolICMP
	case "udp6":
		family, proto = syscall.AF_INET6, IANA_ProtocolIPv6ICMP
	default:
		i := last(network, ':')
		if i < 0 {
			i = len(network)
		}
		switch network[:i] {
		case "ip4":
			proto = IANA_ProtocolICMP
		case "ip6":
			proto = IANA_ProtocolIPv6ICMP
		}
	}

	var cerr error
	var c net.PacketConn
	switch family {
	case syscall.AF_INET, syscall.AF_INET6:
		s, err := syscall.Socket(family, syscall.SOCK_DGRAM, proto)
		if err != nil {
			return nil, os.NewSyscallError("socket", err)
		}
		if (runtime.GOOS == "darwin" || runtime.GOOS == "ios") && family == syscall.AF_INET {
			if err := syscall.SetsockoptInt(s, IANA_ProtocolIP, sysIP_STRIPHDR, 1); err != nil {
				syscall.Close(s)
				return nil, os.NewSyscallError("setsockopt", err)
			}
		}
		sa, err := sockaddr(family, address)
		if err != nil {
			syscall.Close(s)
			return nil, err
		}
		if err := syscall.Bind(s, sa); err != nil {
			syscall.Close(s)
			return nil, os.NewSyscallError("bind", err)
		}
		f := os.NewFile(uintptr(s), "datagram-oriented icmp")
		c, cerr = net.FilePacketConn(f)
		f.Close()
	default:
		c, cerr = net.ListenPacket(network, address)
	}
	if cerr != nil {
		return nil, cerr
	}

	switch proto {
	case IANA_ProtocolICMP:
		return &PacketConn{c: c, p4: ipv4.NewPacketConn(c)}, nil
	case IANA_ProtocolIPv6ICMP:
		return &PacketConn{c: c, p6: ipv6.NewPacketConn(c)}, nil
	default:
		return &PacketConn{c: c}, nil
	}
}

// Get a socket address.
func sockaddr(family int, address string) (syscall.Sockaddr, error) {
	switch family {
	case syscall.AF_INET:
		a, err := net.ResolveIPAddr("ip4", address)
		if err != nil {
			return nil, err
		}
		if len(a.IP) == 0 {
			a.IP = net.IPv4zero
		}
		if a.IP = a.IP.To4(); a.IP == nil {
			return nil, net.InvalidAddrError("non-ipv4 address")
		}
		sa := &syscall.SockaddrInet4{}
		copy(sa.Addr[:], a.IP)
		return sa, nil
	case syscall.AF_INET6:
		a, err := net.ResolveIPAddr("ip6", address)
		if err != nil {
			return nil, err
		}
		if len(a.IP) == 0 {
			a.IP = net.IPv6unspecified
		}
		if a.IP.Equal(net.IPv4zero) {
			a.IP = net.IPv6unspecified
		}
		if a.IP = a.IP.To16(); a.IP == nil || a.IP.To4() != nil {
			return nil, net.InvalidAddrError("non-ipv6 address")
		}
		sa := &syscall.SockaddrInet6{ZoneId: zoneToUint32(a.Zone)}
		copy(sa.Addr[:], a.IP)
		return sa, nil
	default:
		return nil, net.InvalidAddrError("unexpected family")
	}
}

// coming soon.
func zoneToUint32(zone string) uint32 {
	if zone == "" {
		return 0
	}
	if ifi, err := net.InterfaceByName(zone); err == nil {
		return uint32(ifi.Index)
	}
	n, err := strconv.Atoi(zone)
	if err != nil {
		return 0
	}
	return uint32(n)
}

// Get the number of the character given the second argument.
func last(s string, b byte) int {
	i := len(s)
	for i--; i >= 0; i-- {
		if s[i] == b {
			break
		}
	}
	return i
}
