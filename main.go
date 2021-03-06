package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kabukky/httpscerts"
)

func main() {

	if httpscerts.Check("cert.pem", "key.pem") != nil {
		httpscerts.Generate("cert.pem", "key.pem", "localhost")
	}

	var (
		upg = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		mc = MakeMessageChannel()
	)

	http.HandleFunc("/", makeHandler(&upg, &mc))

	if e := http.ListenAndServeTLS(":8182", "cert.pem", "key.pem", nil); e != nil {
		panic(e)
	}
}

func makeHandler(upg *websocket.Upgrader, mc *MessageChannel) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wsc, e := upg.Upgrade(w, r, nil)
		if e != nil {
			fmt.Println(e)
			return
		}

		readMux, writeMux := sync.Mutex{}, sync.Mutex{}
		quitChan := make(chan int)

		fmt.Println("Connection established with", r.RemoteAddr)
		defer fmt.Println("Connection with", r.RemoteAddr, "closed")
		defer func() { quitChan <- 1 }()

		go func() {
			msgChan := make(chan string)
			mc.AddMember(r.RemoteAddr, msgChan)
			defer mc.RemoveMember(r.RemoteAddr)

			for {
				select {
				case m := <-msgChan:
					writeMux.Lock()
					wsc.WriteMessage(websocket.TextMessage, []byte(m))
					writeMux.Unlock()

				case <-quitChan:
					return
				}
			}
		}()

		for {
			readMux.Lock()
			_, data, e := wsc.ReadMessage()
			readMux.Unlock()
			if e != nil {
				fmt.Println(e)
				return
			}
			go mc.SendMessage(fmt.Sprintf("%s", data))
		}
	}
}
