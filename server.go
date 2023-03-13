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

	conn, err := listener.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Accepting a new connection\nFatal error: %s\n", err.Error())
	}

	BenchmarkGOB(conn)
	BenchmarkBinary(conn)
	BenchmarkConn(conn)
}

func BenchmarkGOB(conn net.Conn) {

	dec := gob.NewDecoder(conn)

	buffer := make([]byte, MessageSize)
	for i := 0; i < NbMessages; i++ {

		err := dec.Decode(&buffer)
		if err != nil {
			log.Panic("error receiving message")
		}
	}
}

func BenchmarkBinary(conn net.Conn) {

	buffer := make([]byte, MessageSize)
	for i := 0; i < NbMessages; i++ {

		err := binary.Read(conn, binary.LittleEndian, buffer)
		if err != nil {
			log.Panic("error receiving message")
		}
	}
}

func BenchmarkConn(conn net.Conn) {

	bytes := make([]byte, MessageSize)
	for i := 0; i < NbMessages; i++ {

		ReceiveConn(conn, bytes)
	}
}

func ReceiveConn(conn net.Conn, bytes []byte) {

	pong := 0
	for pong < int(MessageSize) {

		receivedBytes, err := conn.Read(bytes[pong:MessageSize])
		if err != nil || receivedBytes <= 0 {
			log.Panic("error receiving message")
		}
		pong += receivedBytes
	}
}
