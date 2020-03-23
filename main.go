package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

var cnf config
var client Handler

func main() {
	loadConfig(&cnf)
	fmt.Println("iceburg \u2014 classroom chatbot")
	log.Println("creating new connection")
	createTable()
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}

	client = Handler{wac, time.Now()}
	wac.AddHandler(client)
	if err := login(wac); err != nil {
		log.Fatalf("error logging in: %v\n", err)
	}

	pong, err := wac.AdminTest()

	if !pong || err != nil {
		log.Fatalf("error pinging in: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("shutting down now")
	session, err := wac.Disconnect()
	if err != nil {
		log.Fatalf("error disconnecting: %v\n", err)
	}
	if err := writeSession(session); err != nil {
		log.Fatalf("error saving session: %v", err)
	}
}
