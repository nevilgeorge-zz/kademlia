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

	// initialize the data entry table
	k.Table = make(map[ID][]byte)

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

func (k *Kademlia) FindKBucket(nodeId ID) *KBucket {
	fmt.Println("FindKBucket")
	prefixLen := k.NodeID.Xor(nodeId).PrefixLen()
	var index int
	if prefixLen == 160 {
		index = 0
	} else {
		index = 159 - prefixLen
	}

	return &k.BucketList[index]
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
	fmt.Println(len(k.BucketList))
	for i := 0; i < len(k.BucketList); i++ {
		kb := k.BucketList[i]
		for j := 0; j < len(kb.ContactList); j++ {
			c := kb.ContactList[j]
			if c.NodeID.Equals(nodeId) {
				return &c, nil
			}
		}
	}
	err := new(NotFoundError)
	err.msg = "Contact not found!"
	return nil, err
}

// This is the function to perform the RPC
func (k *Kademlia) DoPing(host net.IP, port uint16) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("DoPing")
	address := host.String() + ":" + strconv.Itoa(int(port))
	fmt.Println(address)
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	// create new ping to send to the other node
	ping := new(PingMessage)
	ping.Sender = k.SelfContact
	ping.MsgID = NewRandomID()
	ping.Sender = k.SelfContact
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
	address := string(contact.Host.String()) + ":" + strconv.Itoa(int(contact.Port))
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
	fmt.Println("FindNode")
	address := contact.Host.String() + ":" + strconv.Itoa(int(contact.Port))
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

	if result.Err != nil {
		return "ERR: Error occurred in FindNode RPC"
	}
	// update contact in kbucket of this kademlia
	k.UpdateContactInKBucket(contact)

	// jwhang: taken directly from print_contact in main.go
	// probably a bad idea. need to abstract it out
	response := "OK:\n"
	found := false
	for i := 0; i < len(result.Nodes); i++ {
		c := result.Nodes[i]
		if c.Host != nil {
			found = true
			response += "	Host = " + c.Host.String()
			response += "	Port = " + strconv.Itoa(int(c.Port)) + "\n"
		}
	}

	if found {
		return response
	} else {
		return "ERR: NOT FOUND"
	}
}

func (k *Kademlia) DoFindValue(contact *Contact, searchKey ID) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("DoFindValue")
	address := contact.Host.String() + ":" + strconv.Itoa(int(contact.Port))
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
	if result.Err != nil {
		return "ERR: Error occurred in FindValue RPC"
	}

	k.UpdateContactInKBucket(contact)
	return "OK: Contact updated in KBucket"
}

func (k *Kademlia) LocalFindValue(searchKey ID) string {
	// TODO: Implement
	// If all goes well, return "OK: <output>", otherwise print "ERR: <messsage>"
	fmt.Println("LocalFindValue")
	val := k.Table[searchKey]
	if val == nil || len(val) == 0 {
		return "ERR: Value not found in local table"
	}

	return "OK: " + string(val)
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
	fmt.Println("CALLING UPDATE NOW")
	bucket.Update(*update)
}

// jwhang: updates each contact
func (k *Kademlia) UpdateContacts(contact Contact) {
	fmt.Println("UpdateContacts")
	prefixLen := k.NodeID.Xor(contact.NodeID).PrefixLen()
	if prefixLen == 160 {
		prefixLen = 0
	}

	currentBucket := k.BucketList[prefixLen]
	currentBucket.Update(contact)
}

// nsg622: finds closest nodes
// assumes closest nodes are in the immediate kbucket and the next one
func (k *Kademlia) FindCloseContacts(key ID, req ID) []Contact {
	fmt.Println("FindCloseContacts")
	prefixLen := k.NodeID.Xor(key).PrefixLen()
	fmt.Println("prefixLen is: " + strconv.Itoa(prefixLen))
	var index int
	if prefixLen == 160 {
		index = 0
	} else {
		index = 159 - prefixLen
	}
	contacts := make([]Contact, 20)
	fmt.Println("Index is: " + strconv.Itoa(index))

	for _, val := range k.BucketList[index].ContactList {
		contacts = append(contacts, val)
	}

	if len(contacts) != 20 {
		for _, val := range k.BucketList[index+1].ContactList {
			contacts = append(contacts, val)
			if len(contacts) == 20 {
				break
			}
		}
	}

	return contacts
}

// func (k *Kademlia) AddBucketToSlice(requester ID, bucketNum int, source *[]Contact) {
// 	//k.contacts_mutex[bucketNum].Lock() // jwhang: TODO add mutex for multithread
// 	k.AddBucketContentsToSlice(k.BucketList[bucketNum], requester, source)
// 	//k.contacts_mutex[bucketNum].Unlock()
// }

// func (k *Kademlia) AddBucketContentsToSlice(blist KBucket, requester ID, source *[]Contact) {
// 	// count := 0
// 	// i := 0
// 	// emptySpace := cap(*source) - len(*source)
// 	// for count < emptySpace {
// 	// 	b := blist.ContactList[i]
// 	// 	if b == nil {
// 	// 		break
// 	// 	}
// 	// 	if b.Value.(Contact).NodeID.Equals(requester) {
// 	// 		*source = append(*source, b.Value.(Contact))
// 	// 		count += 1
// 	// 	}
// 	// 	i += 1
// 	// }
// 	// return
// }
