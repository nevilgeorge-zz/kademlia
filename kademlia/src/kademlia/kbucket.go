package kademlia
 
import (
	"fmt"
	//"log"
	//"net"
	//"net/http"
	//"net/rpc"
	//"strconv"
)

type KBucket struct {
	NodeId ID
	NumContacts int
	ContactList [k]Contact
	Kad *Kademlia
	BitMap [k]bool
}

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
