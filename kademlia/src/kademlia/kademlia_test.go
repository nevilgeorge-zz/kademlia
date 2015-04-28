package kademlia

// Kademlia Test Suite
// written by jwhang

// TODO NOTE IMPORTANT
// Haven't figured out what to pass in to NewKademlia()...
// Looks like some address but not entirely sure
import (
	"bytes"
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

// TestStore
// func TestStore(t *testing.T) {
// 	kc := new(KademliaCore)
// 	kc.kademlia = NewKademlia(":7890")
// 	senderID := NewRandomID()
// 	messageID := NewRandomID()
// 	key, err := IDFromString("1234567890")
// 	if err != nil {
// 		t.Error("Couldn't encode key")
// 	}
// 	value := []byte("somedata")
// 	con := Contact{
// 		NodeID: senderID,
// 		Host:   net.IPv4(0x01, 0x02, 0x03, 0x04),
// 		Port:   1234,
// 	}
// 	req := StoreRequest{
// 		Sender: con,
// 		MsgID:  messageID,
// 		Key:    key,
// 		Value:  value,
// 	}
// 	res := new(StoreResult)
// 	err = kc.Store(req, res)
// 	if err != nil {
// 		t.Error("Failed to store key-value pair")
// 	}
// 	if messageID.Equals(res.MsgID) == false {
// 		t.Error("TestStore Failed: MessageID Doesn't match")
// 		t.Fail()
// 	}
// 	if bytes.Equal((*kc).kademlia.Table[key], value) == false {
// 		t.Error("Value stored is incorrect")
// 	}
// }

// TestFindValue
func TestStoreKeyWithFindValue(t *testing.T) {
	kc := new(KademliaCore)
	kc.kademlia = NewKademlia(":7890")
	senderID, messageID := NewRandomID(), NewRandomID()
	key, err := IDFromString("1234567890123456789012345678901234567890")
	if err != nil {
		t.Error("Could not encode key")
		t.Fail()
	}
	value := []byte("somedata")
	con := Contact{
		NodeID: senderID,
		Host:   net.IPv4(0x01, 0x02, 0x03, 0x04),
		Port:   1234,
	}
	req := StoreRequest{
		Sender: con,
		MsgID:  messageID,
		Key:    key,
		Value:  value,
	}
	res := new(StoreResult)
	err = kc.Store(req, res)
	if err != nil {
		t.Error("Failed to store key-value pair")
		t.Fail()
	}
	if messageID.Equals(res.MsgID) == false {
		t.Error("TestStore Failed: MessageID Doesn't match")
		t.Fail()
	}
	messageID = NewRandomID()
	findRequest := FindValueRequest{
		Sender: con,
		MsgID:  messageID,
		Key:    key,
	}
	findResult := new(FindValueResult)
	err = kc.FindValue(findRequest, findResult)
	if err != nil {
		t.Error("Failed to execute find value")
		t.Fail()
	}
	if false == bytes.Equal(findResult.Value, value) {
		t.Error("Retrieved value incorrect")
		t.Fail()
	}
	if messageID.Equals(findResult.MsgID) == false {
		t.Error("TestFindValue Failed: Message ID Doesn't match")
	}
	if len(findResult.Nodes) == 0 {
		t.Error("Returned neighbor nodes without any neighbors! Impossible!")
		t.Fail()
	}
}

// TestFindNode
func TestFindNode(t *testing.T) {
	// TODO
	// Implement
}
