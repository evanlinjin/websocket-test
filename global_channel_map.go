package main

import (
	"fmt"
	"sync"
)

// MessageChannel is a global relay for all messages.
type MessageChannel struct {
	sync.Mutex
	members map[string]chan string
}

// MakeMessageChannel makes a MessageChannel.
func MakeMessageChannel() (m MessageChannel) {
	m.members = make(map[string]chan string)
	return
}

// AddMember adds a member from MessageChannel.
func (m *MessageChannel) AddMember(MemberID string, privateChan chan string) error {
	m.Lock()
	defer m.Unlock()
	if _, got := m.members[MemberID]; got {
		return fmt.Errorf("an entry already exists for '%s'", MemberID)
	}
	m.members[MemberID] = privateChan
	fmt.Printf("[MessageChannel] MEMBER COUNT: %v \n", len(m.members))
	return nil
}

// RemoveMember removes a member from MessageChannel.
func (m *MessageChannel) RemoveMember(MemberID string) {
	m.Lock()
	defer m.Unlock()
	delete(m.members, MemberID)
	fmt.Printf("[MessageChannel] MEMBER COUNT: %v \n", len(m.members))
}

// SendMessage sends message to all members of MessageChannel.
func (m *MessageChannel) SendMessage(msg string) {
	m.Lock()
	defer m.Unlock()
	for _, v := range m.members {
		v <- msg
	}
}
