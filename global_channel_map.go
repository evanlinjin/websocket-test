package main

import (
	"fmt"
	"sync"
)

// GlobalChannelMap is a global relay for all messages.
type GlobalChannelMap struct {
	sync.Mutex
	cMap map[string]chan string
	msgQ chan string
	stop chan bool
}

// MakeGlobalChannelMap makes a GlobalChannelMap.
func MakeGlobalChannelMap() (m GlobalChannelMap) {
	m.cMap = make(map[string]chan string)
	m.stop = make(chan bool)
	return
}

// AddChannel adds a channel to GCM.
func (m *GlobalChannelMap) AddChannel(remoteAdd string, msgChan chan string) error {
	// m.Lock()
	// defer m.Unlock()
	if _, got := m.cMap[remoteAdd]; got {
		return fmt.Errorf("an entry already exists for '%s'", remoteAdd)
	}
	m.cMap[remoteAdd] = msgChan
	fmt.Printf("[GSM] CHANNEL COUNT: %v \n", len(m.cMap))
	return nil
}

// RemoveChannel removes a channel from GCM.
func (m *GlobalChannelMap) RemoveChannel(remoteAdd string) {
	// m.Lock()
	// defer m.Unlock()
	delete(m.cMap, remoteAdd)
	fmt.Printf("[GSM] CHANNEL COUNT: %v \n", len(m.cMap))
}

// SendMessage sends message to all channels.
func (m *GlobalChannelMap) SendMessage(msg string) {
	m.Lock()
	defer m.Unlock()
	for _, v := range m.cMap {
		v <- msg
	}
}
