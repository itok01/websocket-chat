package main

import (
	"flag"
)

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func main() {
	addr := flag.String("addr", "", "access server")
	flag.Parse()
	if *addr != "" {
		client(*addr)
	} else {
		server()
	}
}
