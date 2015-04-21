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
	NodeId ID
	NumContacts int
	ContactList [k]Contact
	Kad *Kademlia
	BitMap [k]bool
}

// Initialize KBuckets, called 160 times when Kademlia is instantiated in kademlia.org
func (kb *KBucket) Initialize() {
	fmt.Println("Initialize")
	kb.NumContacts = 0
	kb.NodeId = NewRandomID()
	for i := 0; i < 20; i++ {
		nullContact := new(Contact)
		kb.ContactList[i] = *nullContact
		kb.BitMap[i] = false
	}
}

// Remove the contact corresponding to a given ID from the KBucket
func (kb *KBucket) RemoveContact(targetID ID) (bool) {
	fmt.Println("RemoveContact")
	for i,_ := range kb.ContactList {
		if kb.ContactList[i].NodeID == targetID && kb.BitMap[i] == true  {
			nullContact := new(Contact)
			kb.ContactList[i] = *nullContact
			kb.BitMap[i] = false
			return true
		}
	}
	return false
}

// Adds a given contact to a free space in the KBucket
func (kb *KBucket) AddContact(newContact Contact) {
	fmt.Println("AddContact")
	for i,_ := range kb.ContactList {
		if kb.BitMap[i] == false {
			kb.ContactList[i] = newContact
			kb.BitMap[i] = true
			break
		}
	}
}

