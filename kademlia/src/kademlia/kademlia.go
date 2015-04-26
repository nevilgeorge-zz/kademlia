package kademlia

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

const (
	alpha   = 3
	b       = 8 * IDBytes
	k       = 20
	kb_size = 160
)

// Kademlia type. You can put whatever state you need in this.
type Kademlia struct {
	NodeID        ID
	SelfContact   Contact
	BucketList    []KBucket
	Table         map[ID][]byte
}

func NewKademlia(laddr string) *Kademlia {
	// TODO: Initialize other state here as you add functionality.
	fmt.Println("NewKademlia")
	k := new(Kademlia)
	k.NodeID = NewRandomID()
	// only 160 nodes in this system
	k.BucketList = make([]KBucket, kb_size)

	// initialize all k-buckets
	for i := 0; i < b; i++ {
		k.BucketList[i].Initialize()
	}

	// Set up RPC server
	// NOTE: KademliaCore is just a wrapper around Kademlia. This type includes
	// the RPC functions.
	rpc.Register(&KademliaCore{k})
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal("Listen: ", err)
	}
	// Run RPC server forever.
	go http.Serve(l, nil)

	// Add self contact
	hostname, port, _ := net.SplitHostPort(l.Addr().String())
	port_int, _ := strconv.Atoi(port)
	ipAddrStrings, err := net.LookupHost(hostname)
	var host net.IP
	for i := 0; i < len(ipAddrStrings); i++ {
		host = net.ParseIP(ipAddrStrings[i])
		if host.To4() != nil {
			break
		}
	}
	k.SelfContact = Contact{k.NodeID, host, uint16(port_int)}
	return k
}

func (k *Kademlia) FindKBucket(nodeId ID) KBucket {
	fmt.Println("FindKBucket")
	distance := k.NodeID.Xor(nodeId)
	index := distance.PrefixLen()
	return k.BucketList[index]
}

type NotFoundError struct {
	id  ID
	msg string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%x %s", e.id, e.msg)
}

func (k *Kademlia) FindContact(nodeId ID) (*Contact, error) {
	// TODO: Search through contacts, find specified ID
	// Find contact with provided ID
	fmt.Println("FindContact")
	if nodeId == k.NodeID {
		return &k.SelfContact, nil
	}
	kb := k.FindKBucket(nodeId)
	return &kb.ContactList[0], nil
	//return nil, &NotFoundError{nodeId, "Not found"}
}

// This is the function to perform the RPC
func (k *Kademlia) DoPing(host net.IP, port uint16) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("DoPing")
	address := string(host) + ":" + strconv.Itoa(int(port))
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("DialHTTP in DoPing: ", err)
	}

	// create new ping to send to the other node
	ping := new(PingMessage)
	ping.MsgID = NewRandomID()
	var pong PongMessage
	err = client.Call("KademliaCore.Ping", ping, &pong)
	if err != nil {
		log.Fatal("Call in DoPing ", err)
	}

	// update contact in kbucket of this kademlia
	k.UpdateContactInKBucket(&pong.Sender)

	return "ERR: Not implemented"
}

func (k *Kademlia) DoStore(contact *Contact, key ID, value []byte) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("DoStore")
	address := string(contact.Host) + ":" + strconv.Itoa(int(contact.Port))
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("DialHTTP in DoStore: ", err)
	}

	request := new(StoreRequest)
	request.Sender = *contact
	request.Key = key
	request.Value = value
	request.MsgID = NewRandomID()

	var result StoreResult
	err = client.Call("KademliaCore.Store", request, &result)
	if err != nil {
		log.Fatal("Call in DoStore ", err)
	}

	// update contact in kbucket of this kademlia
	k.UpdateContactInKBucket(contact)

	return "ERR: Not implemented"
}

func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("FindFindNode")
	address := string(contact.Host) + ":" + strconv.Itoa(int(contact.Port))
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("DialHTTP in DoStore: ", err)
	}

	request := new(FindNodeRequest)
	request.Sender = *contact
	request.NodeID = searchKey
	request.MsgID = NewRandomID()

	var result FindNodeResult
	err = client.Call("KademliaCore.FindNode", request, &result)
	if err != nil {
		log.Fatal("Call in DoStore ", err)
	}

	// update contact in kbucket of this kademlia
	k.UpdateContactInKBucket(contact)

	return "ERR: Not implemented"
}

func (k *Kademlia) DoFindValue(contact *Contact, searchKey ID) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("DoFindValue")
	return "ERR: Not implemented"
}

func (k *Kademlia) LocalFindValue(searchKey ID) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("LocalFindValue")
	return "ERR: Not implemented"
}

func (k *Kademlia) DoIterativeFindNode(id ID) string {
	// For project 2!
	return "ERR: Not implemented"
}
func (k *Kademlia) DoIterativeStore(key ID, value []byte) string {
	// For project 2!
	return "ERR: Not implemented"
}
func (k *Kademlia) DoIterativeFindValue(key ID) string {
	// For project 2!
	return "ERR: Not implemented"
}

func (k *Kademlia) UpdateContactInKBucket(update *Contact) {
	bucket := k.FindKBucket(update.NodeID)
	bucket.Update(*update)
}