package main

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

var client_connection net.Conn
var dec *gob.Decoder

func run_server() {

	completeAddress := ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", completeAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving address %s\nFatal error: %s\n", completeAddress, err.Error())
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on port %d\nFatal error: %s\n", port, err.Error())
	}

	new_connection, err := listener.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Accepting a new connection\nFatal error: %s\n", err.Error())
	}
	client_connection = new_connection
	dec = gob.NewDecoder(client_connection)

	var fns = []func() error{ReadGOB, ReadBinary, ReadConn}

	fmt.Println("Throughput: ")

	for _, fn := range fns {

		bytes = make([]byte, MessageSize)

		for i := 0; i < NbMessages; i++ {
			fn()
			if err != nil {
				log.Panic("error receiving message")
			}
		}
		bytes = nil
	}
}

func ReadGOB() error {

	return dec.Decode(&bytes)
}

func ReadBinary() error {

	return binary.Read(client_connection, binary.LittleEndian, bytes)
}

func ReadConn() error {

	return ReceiveConn(bytes)
}

func ReceiveConn(bytes []byte) error {

	pong := 0
	for pong < int(MessageSize) {

		receivedBytes, err := client_connection.Read(bytes[pong:MessageSize])
		if err != nil || receivedBytes <= 0 {
			return fmt.Errorf("error receiving message")
		}
		pong += receivedBytes
	}
	return nil
}
