package kademlia

import (
	"fmt"
)

// KBucket struct
type KBucket struct {
	ContactList []Contact
	ContactChan	chan *Contact
	IDChan		chan ID
	Kad         *Kademlia
}

// Initialize KBuckets, called 160 times when Kademlia is instantiated in kademlia.go
func (kb *KBucket) Initialize() {
	// create slice for ContactList
	kb.ContactList = make([]Contact, 0, k)
	kb.ContactChan = make(chan *Contact)
	kb.IDChan = make(chan ID)

	go kb.handleContacts()
}

// go routine function to handle interacting with the ContactList
func (kb *KBucket) handleContacts() {
	for {
		select {
		case newContact := <- kb.ContactChan:
			kb.ContactList = append(kb.ContactList, *newContact)
			fmt.Println("# of contacts in node: ", len(kb.ContactList))

		case targetID := <- kb.IDChan:
			for i, _ := range kb.ContactList {
				if kb.ContactList[i].NodeID == targetID {
					temp := kb.ContactList
					a := append(temp[:i], temp[(i+1):]...)
					kb.ContactList = a
				}
			}
			
		}
	}
}

// Remove the contact corresponding to a given ID from the KBucket
func (kb *KBucket) RemoveContact(targetID ID) {
	fmt.Println("RemoveContact")
	kb.IDChan <- targetID
}

// Adds a given contact to the end of the kbucket
func (kb *KBucket) AddContact(newContact Contact) {
	fmt.Println("AddContact")
	toAdd := new(Contact)
	toAdd.NodeID = newContact.NodeID
	toAdd.Host = newContact.Host
	toAdd.Port = newContact.Port
	kb.ContactChan <- toAdd
}

// returns a boolean for whether a given Contact exists in the KBucket and index if it was found
func (kb *KBucket) ContainsContact(cont Contact) (exists bool, index int) {
	// iterate through ContactList and compare Contact NodeIDs
	for i := 0; i < len(kb.ContactList); i++ {
		current := kb.ContactList[i]
		if current.NodeID.Equals(cont.NodeID) {
			exists = true
			index = i
			return
		}
	}
	exists = false
	index = -1
	return
}

// Update the KBucket to sort the nodes with most recently used in at the head of the KBucket
func (kb *KBucket) Update(updated Contact) {
	fmt.Println("Update")
	fmt.Println(len(kb.ContactList))
	fmt.Print("K is : ")
	fmt.Print(k)

	// check whether the updated contact exists in the KBucket
	exists, _ := kb.ContainsContact(updated)
	if exists {
		fmt.Println("It exists!")
		// move Contact to the end of the KBucket
		kb.MoveToTail(updated)
	} else if len(kb.ContactList) < k {
		fmt.Println("New contact")
		// create a new contact for the node and add it to the tail of the KBucket
		// not sure if a new Contact needs to be created, but that's what the doc says
		//temp := Contact(CopyID(updated.NodeID), updated.Host, updated.Port)
		temp := new(Contact)
		temp.NodeID = CopyID(updated.NodeID)
		temp.Host = updated.Host
		temp.Port = updated.Port
		fmt.Println("NodeID:")
		fmt.Println(temp.NodeID)
		fmt.Println("Host:")
		fmt.Println(temp.Host)
		fmt.Println("Port:")
		fmt.Println(temp.Port)
		kb.AddContact(*temp) // jwhang: kinda fishy.. not sure if this is ok
	} else {
		// ping first node in slice
		// if it doesn't respond, removeContact(oldContact) and addContact(updated)
		// else moveToTail(oldContact) and ignore updated
		firstContact := kb.ContactList[0]
		ret := kb.Kad.DoPing(firstContact.Host, firstContact.Port)
		if ret == "Error" { // jwhang TODO: Fix this to nil
			kb.RemoveContact(firstContact.NodeID)
			kb.AddContact(updated)
		} else {
			kb.MoveToTail(firstContact)
		}
	}
}

// moves a contact from its position in the KBucket to the end of the same KBucket
func (kb *KBucket) MoveToTail(updated Contact) {
	exists, _ := kb.ContainsContact(updated)
	if exists {
		// finds and removes contact if already exists
		kb.RemoveContact(updated.NodeID)
	}
	// adds to the end of the KBucket
	kb.AddContact(updated)
}
