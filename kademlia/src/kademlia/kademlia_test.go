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
	// "strings"
	"testing"
	// "fmt"
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

// func TestStore(t *testing.T) {
// 	kc := new(KademliaCore)
// 	kc.kademlia = NewKademlia("localhost:1234")
// 	senderID := NewRandomID()
// 	messageID := NewRandomID()
// 	key, err := IDFromString("1234567890123456789012345678901234567890")
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
// 		t.Error("TestStore: Failed to store key-value pair")
// 		t.Fail()
// 	}
// 	if messageID.Equals(res.MsgID) == false {
// 		t.Error("TestStore: MessageID Doesn't match")
// 		t.Fail()
// 	}
// 	if bytes.Equal((*kc).kademlia.Table[key], value) == false {
// 		t.Error("TestStore: Value stored is incorrect")
// 		t.Fail()
// 	}
// }

// TestFindValue
func TestStoreKeyWithFindValue(t *testing.T) {
	kc := new(KademliaCore)
	kc.kademlia = NewKademlia("localhost:1235")
	senderID, messageID := NewRandomID(), NewRandomID()
	key, err := IDFromString("1234567890123456789012345678901234567890")
	if err != nil {
		t.Error("Could not encode key")
		t.Fail()
	}
	value := []byte("somedata")
	con := Contact{
		NodeID: senderID,
		Host:   net.IPv4(127, 0, 0, 1),
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
	if len(findResult.Nodes) == 1 {
		t.Error("Returned neighbor nodes without any neighbors! Impossible!")
		t.Fail()
	}
}

// TestPingSelf
// Pings itself and sees if it exists in the contact
// func TestPingSelf(t *testing.T) {
// 	kc := new(KademliaCore)
// 	kc.kademlia = NewKademlia("localhost:1234")
// 	senderID := NewRandomID()
// 	messageID := NewRandomID()
// 	key, err := IDFromString("1234567890123456789012345678901234567890")
// 	if err != nil {
// 		t.Error("Couldn't encode key")
// 	}
// 	value := []byte("somedata")
// 	selfHost := net.IPv4(127, 0, 0, 1)
// 	selfPort := uint16(1234)
// 	res := kc.kademlia.DoPing(selfHost, selfPort)
// 	if strings.Contains(res, "ERR") {
// 		t.Error("TestPingSelf: Failed to ping itself")
// 		t.Fail()
// 	}
// }

// func TestPingAnother(t *testing.T) {
// 	kc1 := new(KademliaCore)
// 	kc2 := new(KademliaCore)
// 	kc1.kademlia = NewKademlia("localhost:7980")
// 	kc2.kademlia = NewKademlia("localhost:1234")
// 	k21ID := kc1.kademlia.NodeID
// 	kc1Host := net.IPv4(127, 0, 0, 1)
// 	kc2Port := uint16(7890)
// 	kc2ID := kc2.kademlia.NodeID
// 	kc2Host := net.IPv4(127, 0, 0, 1)
// 	kc2Port := uint16(1234)
// 	res := kc1.kademlia.DoPing(kc2Host, kc2Port)
// 	if strings.Contains(res, "ERR") {
// 		t.Error("TestPingAnother: Failed to ping node 2 from node 1")
// 		t.Fail()
// 	}
// 	res, err = kc1.FindContact(kc2ID)
// 	if err != nil {
// 		t.Error("TestPingAnother: The sender doesn't have the receiver's contact!")
// 		t.Fail()
// 	}
// 	if !(res.NodeID.Equals(kc2ID)) || res.Host.String() != kc2Host.String() || res.Port != kc2Port {
// 		t.Error("TestPingAnother: The sender's contact of receiver doesn't match with the actual receiver info!")
// 		t.Fail()
// 	}
// 	res, err = kc2.FindContact(kc1ID)
// 	if err != nil {
// 		t.Error("TestPingAnother: The receiver doesn't have sender's contact!")
// 		t.Fail()
// 	}
// 	if !(res.NodeID.Equals(kc1ID)) || res.Host.String() != kc1Host.String() || res.Port != kc1Port {
// 		t.Error("TestPingAnother: The receiver's contact of the sender doesn't match with the actual sender info")
// 		t.Fail()
// 	}
// }

// func TestFindNode(t *testing.T) {
// 	kc1 := new(KademliaCore)
// 	kc2 := new(KademliaCore)
// 	kc1.kademlia = NewKademlia("localhost:7980")
// 	kc2.kademlia = NewKademlia("localhost:1234")
// 	k21ID := kc1.kademlia.NodeID
// 	kc1Host := net.IPv4(127, 0, 0, 1)
// 	kc2Port := uint16(7890)
// 	kc2ID := kc2.kademlia.NodeID
// 	kc2Host := net.IPv4(127, 0, 0, 1)
// 	kc2Port := uint16(1234)

// 	res := kc1.kademlia.DoPing(kc2Host, kc2Port)
// 	if strings.Contains(res, "ERR") {
// 		t.Error("TestPingAnother: Failed to ping node 2 from node 1")
// 		t.Fail()
// 	}

// 	senderID := NewRandomID()
// 	messageID := NewRandomID()
// 	key, err := IDFromString("1234567890123456789012345678901234567890")
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
// 	err = kc1.Store(req, res)
// 	if err != nil {
// 		t.Error("Failed to store key-value pair")
// 		t.Fail()
// 	}

// 	con = Contact{
// 		NodeID: kc1ID,
// 		Host:   kc1Host,
// 		Port:   kc1Port,
// 	}
// 	res = kc2.kademlia(DoFindNode, con, key)
// 	if strings.Contains(res, "ERR") {
// 		t.Error("DoFindNode failed")
// 		t.Fail()
// 	}
// }

// func TestFindValue(t *testing.T) {
// 	kc1 := new(KademliaCore)
// 	kc2 := new(KademliaCore)
// 	kc1.kademlia = NewKademlia("localhost:7980")
// 	kc2.kademlia = NewKademlia("localhost:1234")
// 	k21ID := kc1.kademlia.NodeID
// 	kc1Host := net.IPv4(127, 0, 0, 1)
// 	kc2Port := uint16(7890)
// 	kc2ID := kc2.kademlia.NodeID
// 	kc2Host := net.IPv4(127, 0, 0, 1)
// 	kc2Port := uint16(1234)

// 	res := kc1.kademlia.DoPing(kc2Host, kc2Port)
// 	if strings.Contains(res, "ERR") {
// 		t.Error("TestPingAnother: Failed to ping node 2 from node 1")
// 		t.Fail()
// 	}

// 	senderID := NewRandomID()
// 	messageID := NewRandomID()
// 	key, err := IDFromString("1234567890123456789012345678901234567890")
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
// 	err = kc1.Store(req, res)
// 	if err != nil {
// 		t.Error("Failed to store key-value pair")
// 		t.Fail()
// 	}

// 	con = Contact{
// 		NodeID: kc1ID,
// 		Host:   kc1Host,
// 		Port:   kc1Port,
// 	}
// 	res = kc2.kademlia(DoFindValue, con, key)
// 	if strings.Contains(res, "ERR") {
// 		t.Error("DoFindNode failed")
// 		t.Fail()
// 	}
// }
