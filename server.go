package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const portNumber = ":3000"

var (
	conns      []net.Conn
	connCh     = make(chan net.Conn)
	closeCh    = make(chan net.Conn)
	msgCh      = make(chan string)
	clientName string
)

func main() {
	server, err := net.Listen("tcp", portNumber)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Fatal(err)
			}

			conns = append(conns, conn)
			connCh <- conn
		}
	}()

	for {
		select {
		case conn := <-connCh:
			go onMessage(conn)

		case msg := <-msgCh:
			fmt.Print(msg)

		case conn := <-closeCh:
			fmt.Println("Client exit")
			removeConn(conn)
		}
	}

	// end main
}

func publishMsg(conn net.Conn, msg string) {
	for i, _ := range conns {
		if conns[i] != conn {
			conns[i].Write([]byte(msg))
		}
	}
}

func removeConn(conn net.Conn) {
	var i int
	for i, _ = range conns {
		if conns[i] == conn {
			break
		}
	}

	conns = append(conns[:i], conns[i+1:]...)
}

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msgCh <- msg
		publishMsg(conn, msg)
	}
	closeCh <- conn
}
