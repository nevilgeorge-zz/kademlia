package kademlia

// Kademlia Test Suite
// written by jwhang

import (
	"net"
	"os"
	"testing"
)

func getHostIP() net.IP {
	host, err := os.Hostname()
	if err != nil {
		return net.IPv4(byte(127), 0, 0, 1)
	}
	addr, err := net.LookupAddr(host)
	if len(addr) < 0 || err != nil {
		return net.IPv4(byte(127), 0, 0, 1)
	}
	return net.ParseIP(addr[0])
}

var hostIP = getHostIP()

// TestPing
func TestLocalFindValue(t *testing.T) {
	for i := 0; i < 100; i++ {
		_in, _ := IDFromInteger(i)

	}
}

// TestStore
func TestStore(t *testing.T) {
	k := NewKademlia()
	sender := NewRandomID()
	receiver := NewRandomID()

}
