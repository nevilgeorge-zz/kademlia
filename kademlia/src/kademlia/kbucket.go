package kademlia
 
import (
	//"fmt"
	//"log"
	//"net"
	//"net/http"
	//"net/rpc"
	//"strconv"
)

type KBucket struct {
	NumContacts int
	ContactList [20]Contact
	Kad *Kademlia
	BitMap [20]bool
}

func (kb *KBucket) Initialize() {
	for i := 0; i < 20; i++ {
		nullContact := new(Contact)
		kb.ContactList[i] = *nullContact
		kb.BitMap[i] = false
	}
}

func (kb *KBucket) RemoveContact(targetID ID) (bool) {
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

func (kb *KBucket) AddContact(newContact Contact) {
	for i,_ := range kb.ContactList {
		if kb.BitMap[i] == false {
			kb.ContactList[i] = newContact
			kb.BitMap[i] = true
			break
		}
	}
}
