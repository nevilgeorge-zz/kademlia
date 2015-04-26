package main

import (
	"testing"
)

func (k *Kademlia) TestPing(destIP net.IP, destPort uint16) {
	for i := 0; i < 100; i++ {
		pingRes := k.DoPing(netIP, destPort)
	}
}
