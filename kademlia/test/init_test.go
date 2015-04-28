package main

import (
	"testing"
	"fmt"
)

import (
	"kademlia"
)

func TestInit(t *testing.T) {
	kad1 := kademlia.NewKademlia(":7890")
	kad2 := kademlia.NewKademlia(":7890")

	fmt.Println(kad1.NodeID)
	fmt.Println(kad2.NodeID)
}