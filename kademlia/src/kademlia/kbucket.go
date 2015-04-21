package kademlia

import (
	"fmt"
	//"log"
	//"net"
	//"net/http"
	//"net/rpc"
	//"strconv"
)

// KBucket struct
type KBucket struct {
	NodeID      ID
	ContactList []Contact
	Kad         *Kademlia
	Kcore       KademliaCore
}

// Initialize KBuckets, called 160 times when Kademlia is instantiated in kademlia.go
func (kb *KBucket) Initialize() {
	fmt.Println("Initialize")
	kb.NodeID = NewRandomID()
	// create slice for ContactList
	kb.ContactList = make([]Contact, k)
}

// Remove the contact corresponding to a given ID from the KBucket
func (kb *KBucket) RemoveContact(targetID ID) bool {
	fmt.Println("RemoveContact")
	for i, _ := range kb.ContactList {
		if kb.ContactList[i].NodeID == targetID {
			temp := kb.ContactList
			a := append(temp[:i], temp[(i+1):]...)
			kb.ContactList = a
			return true
		}
	}
	return false
}

// Adds a given contact to the end of the kbucket
func (kb *KBucket) AddContact(newContact Contact) {
	fmt.Println("AddContact")
	kb.ContactList = append(kb.ContactList, newContact)
}

// returns a boolean for whether a given Contact exists in the KBucket and index if it was found
func (kb *KBucket) ContainsContact(cont Contact) (exists bool, index int) {
	// iterate through ContactList and compare Contact NodeIDs
	for i := 0; i < len(kb.ContactList); i++ {
		current := kb.ContactList[i]
		if current.NodeID.Equals(cont.NodeID) {
			exists = true
			index = i
			return exists
		}
	}
	exists = false
	index = -1
	return exists
}

// Update the KBucket to sort the nodes with most recently used in at the head of the KBucket
func (kb *KBucket) Update(updated Contact) {
	// check whether the updated contact exists in the KBucket
	exists, _ := ContainsContact(updated)
	if exists {
		// move Contact to the end of the KBucket
		kb.moveToTail(updated)
	} else if len(kb.ContactList) < k {
		// create a new contact for the node and add it to the tail of the KBucket
		// not sure if a new Contact needs to be created, but that's what the doc says
		//temp := Contact(CopyID(updated.NodeID), updated.Host, updated.Port)
		temp := new(Contact)
		temp.NodeID = CopyID(updated.NodeID)
		temp.Host = updated.Host
		temp.Port = updated.Port
		kb.AddContact(*temp) // jwhang: kinda fishy.. not sure if this is ok
	} else {
		// ping first node in slice
		// if it doesn't respond, removeContact(oldContact) and addContact(updated)
		// else moveToTail(oldContact) and ignore updated
		firstContact := kb.ContactList[0]
		ret := kb.Kad.DoPing(firstContact.Host, firstContact.Port)
		if ret == null {
			kb.RemoveContact(firstContact)
			kb.AddContact(updated)
		} else {
			kb.MoveToTail(firstContact)
		}
	}
}

// moves a contact from its position in the KBucket to the end of the same KBucket
func (kb *KBucket) MoveToTail(updated Contact) {
	exists, _ := ContainsContact(updated)
	if exists {
		// finds and removes contact if already exists
		kb.RemoveContact(updated)
	}
	// adds to the end of the KBucket
	kb.AddContact(updated)
}
