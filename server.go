package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"time"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	fmt.Println(err)
	time.Sleep(1 * time.Second)
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			continue
		}
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	defer client.Close()

	reader := bufio.NewReader(client)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if len(msg) > 0 {
			msgs <- Message{sender: clientid, message: msg}
		}
	}
}

func main() {
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	ln, _ := net.Listen("tcp", *portPtr)
	counter := 0
	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			fmt.Println("successfully! Device counter: ", counter+1)
			clients[counter] = conn
			go handleClient(clients[counter], counter, msgs)
			counter++

		case msg := <-msgs:
			fmt.Println("received from ", msg.sender, ": ", msg.message)
			for id, conn := range clients {
				if id != msg.sender {
					fmt.Fprint(conn, msg.message)
				}
			}

		}
	}
}
