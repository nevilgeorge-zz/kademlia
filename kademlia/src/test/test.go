package main

import (
	"testing"
)

func SetUp() {
	kad1 = NewKademlia("first")
	kad2 = NewKademlia("second")

	kad1.DoPing(kad2.SelfContact.Host, kad1.SelfContact.Port)
	kad2.DoPing(kad1.SelfContact.Host, kad1.SelfContact.Port)

	kad1.DoStore(kad2.SelfContact, kad2.NodeID, 10000)
	kad1.DoFindValue(kad2.SelfContact, kad2.NodeID)
}

func TestPing()
