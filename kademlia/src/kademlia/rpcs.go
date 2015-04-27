package kademlia

// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"net"
)

type KademliaCore struct {
	kademlia *Kademlia
}

// Host identification.
type Contact struct {
	NodeID ID
	Host   net.IP
	Port   uint16
}

///////////////////////////////////////////////////////////////////////////////
// PING
///////////////////////////////////////////////////////////////////////////////
type PingMessage struct {
	Sender Contact
	MsgID  ID
}

type PongMessage struct {
	MsgID  ID
	Sender Contact
}

func (kc *KademliaCore) Ping(ping PingMessage, pong *PongMessage) error {
	// Specify the sender
	// Update contact, etc
	// sender is this node
	c := kc.kademlia.SelfContact
	pong.Sender = c
	pong.MsgID = CopyID(ping.MsgID)

	// update contact in this kademlia kbucket
	kc.kademlia.UpdateContactInKBucket(&ping.Sender)

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// STORE
///////////////////////////////////////////////////////////////////////////////
type StoreRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
	Value  []byte
}

type StoreResult struct {
	MsgID ID
	Err   error
}

func (kc *KademliaCore) Store(req StoreRequest, res *StoreResult) error {
	kc.kademlia.Table[req.Key] = req.Value
	res.MsgID = CopyID(req.MsgID)
	res.Err = nil

	// update contact in kbucket
	kc.kademlia.UpdateContactInKBucket(&req.Sender)

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// FIND_NODE
///////////////////////////////////////////////////////////////////////////////
type FindNodeRequest struct {
	Sender Contact
	MsgID  ID
	NodeID ID
}

type FindNodeResult struct {
	MsgID ID
	Nodes []Contact
	Err   error
}

func (kc *KademliaCore) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	// TODO: Implement.
	kc.kademlia.UpdateContacts(req.Sender)
	res.MsgID = CopyID(req.MsgID)
	res.Nodes = kc.kademlia.FindCloseContacts(req.NodeID, kc.kademlia.NodeID)
	res.Err = nil

	// update contact in kbucket
	kc.kademlia.UpdateContactInKBucket(&req.Sender)

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// FIND_VALUE
///////////////////////////////////////////////////////////////////////////////
type FindValueRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
	MsgID ID
	Value []byte
	Nodes []Contact
	Err   error
}

func (kc *KademliaCore) FindValue(req FindValueRequest, res *FindValueResult) error {
	// TODO: Implement.
	res.MsgID = CopyID(req.MsgID)
	res.Nodes = kc.kademlia.FindCloseContacts(req.Sender.NodeID, kc.kademlia.NodeID)
	val := kc.kademlia.Table[req.Key]
	if val == nil || len(val) == 0 {
		err := new(NotFoundError)
		err.msg = "Value is nil or is empty byte array"
		res.Err = err
		return err
	}
	res.Value = val
	res.Err = nil

	// update contact in kbucket
	kc.kademlia.UpdateContactInKBucket(&req.Sender)

	return nil
}
