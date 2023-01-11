package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func onMessageClient(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter your name: ")
	nameReader := bufio.NewReader(os.Stdin)
	nameInput, _ := nameReader.ReadString('\n')

	nameInput = nameInput[:len(nameInput)-1]

	fmt.Println("Your name is: ", nameInput)

	fmt.Println("----------------------")

	go onMessageClient(conn)
	for {
		fmt.Print("Enter your msg: ")
		msgReader := bufio.NewReader(os.Stdin)
		msg, err := msgReader.ReadString('\n')
		if err != nil {
			break
		}

		msg = fmt.Sprintf("%s: %s\n", nameInput,
			msg[:len(msg)-1])

		conn.Write([]byte(msg))

	}

}
