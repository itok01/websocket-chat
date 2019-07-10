package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

func client(addr string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	myMsg := Message{}

	fmt.Printf("Input name: ")
	myMsg.Username = input()

	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   "/ws",
	}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		for {
			myMsg.Text = input()
			if myMsg.Text != "" {
				err := conn.WriteJSON(myMsg)
				if err != nil {
					log.Printf("error: %v", err)
					conn.Close()
				}
			}
		}
	}()

	myMsg.Text = fmt.Sprintf("Connected %s", myMsg.Username)
	if err := conn.WriteJSON(myMsg); err != nil {
		log.Printf("error: %v", err)
		conn.Close()
	}

	go func() {
		for {
			msg := Message{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				break
			}
			log.Printf("%s: %s", msg.Username, msg.Text)
		}
	}()

	<-interrupt
	myMsg.Text = fmt.Sprintf("Disconnected %s", myMsg.Username)
	if err := conn.WriteJSON(myMsg); err != nil {
		log.Printf("error: %v", err)
	}
	conn.Close()
	return
}

func input() (text string) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	text = scanner.Text()

	text = strings.TrimSpace(text)
	return
}
