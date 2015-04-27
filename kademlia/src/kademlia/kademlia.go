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
	alpha        = 3
	b            = 8 * IDBytes
	k            = 20
	bucket_count = 160
)

// Kademlia type. You can put whatever state you need in this.
type Kademlia struct {
	NodeID      ID
	SelfContact Contact
	BucketList  []KBucket
	Table       map[ID][]byte
}

func NewKademlia(laddr string) *Kademlia {
	// TODO: Initialize other state here as you add functionality.
	fmt.Println("NewKademlia")
	k := new(Kademlia)
	k.NodeID = NewRandomID()
	// only 160 nodes in this system
	k.BucketList = make([]KBucket, bucket_count)

	// initialize all k-buckets
	for i := 0; i < b; i++ {
		k.BucketList[i].Initialize()
	}

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
		log.Fatal("ERR: ", err)
	}

	// create new ping to send to the other node
	ping := new(PingMessage)
	ping.MsgID = NewRandomID()
	var pong PongMessage
	err = client.Call("KademliaCore.Ping", ping, &pong)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	// update contact in kbucket of this kademlia
	updated := pong.Sender
	k.UpdateContactInKBucket(&updated)
	// find kbucket that should hold this contact

	return "OK: Contact updated in KBucket"
}

func (k *Kademlia) DoStore(contact *Contact, key ID, value []byte) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("DoStore")
	address := string(contact.Host) + ":" + strconv.Itoa(int(contact.Port))
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	request := new(StoreRequest)
	request.Sender = *contact
	request.Key = key
	request.Value = value
	request.MsgID = NewRandomID()

	var result StoreResult
	err = client.Call("KademliaCore.Store", request, &result)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	// update contact in kbucket of this kademlia
	k.UpdateContactInKBucket(contact)

	return "OK: Contact updated in KBucket"
}

func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("FindFindNode")
	address := string(contact.Host) + ":" + strconv.Itoa(int(contact.Port))
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	request := new(FindNodeRequest)
	request.Sender = *contact
	request.NodeID = searchKey
	request.MsgID = NewRandomID()

	var result FindNodeResult
	err = client.Call("KademliaCore.FindNode", request, &result)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	// update contact in kbucket of this kademlia
	k.UpdateContactInKBucket(contact)

	return "OK: Contact updated in KBucket"
}

func (k *Kademlia) DoFindValue(contact *Contact, searchKey ID) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("DoFindValue")
	address := string(contact.Host) + ":" + strconv.Itoa(int(contact.Port))
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	request := new(FindValueRequest)
	request.Sender = *contact
	request.Key = searchKey
	request.MsgID = NewRandomID()

	var result FindValueResult
	err = client.Call("KademliaCore.FindValue", request, &result)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	// update contact in kbucket of this kademlia
	k.UpdateContactInKBucket(contact)

	return "OK: Contact updated in KBucket"
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

// jwhang: updates each contact
func (k *Kademlia) UpdateContacts(contact Contact) {
	prefixLen := k.NodeID.Xor(contact.NodeID).PrefixLen()
	currentBucket := k.BucketList[prefixLen]
	currentBucket.Update(contact)
}

// jwhang: FindCloseNodes: to be used in FindNode
func (k *Kademlia) FindCloseContacts(key ID, req ID) []Contact {
	prefixLen := k.NodeID.Xor(key).PrefixLen()
	contacts := make([]Contact, 0, 160)

	if bucket_count > prefixLen {
		k.AddBucketToSlice(req, prefixLen, &contacts)
	} else {
	}
	return contacts
}

func (k *Kademlia) AddBucketToSlice(requester ID, bucketNum int, source *[]Contact) {
	//k.contacts_mutex[bucketNum].Lock() // jwhang: TODO add mutex for multithread
	k.AddBucketContentsToSlice(k.BucketList[bucketNum], requester, source)
	//k.contacts_mutex[bucketNum].Unlock()
}

func (k *Kademlia) AddBucketContentsToSlice(blist KBucket, requester ID, source *[]Contact) {
	count := 0
	i := 0
	emptySpace := cap(*source) - len(*source)
	for count < emptySpace {
		b := blist.ContactList[i]
		if b == nil {
			break
		}
		if b.Value.(Contact).NodeID.Equals(requester) {
			*source = append(*source, b.Value.(Contact))
			count += 1
		}
		i += 1
	}
	return
}
