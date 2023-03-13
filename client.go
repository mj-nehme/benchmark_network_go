package main

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

func run_client() {

	completeAddress := IP + ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", completeAddress)
	if err != nil {
		log.Panicf("Error resolving address %s", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Panicf("DialTCP Error %s", err)
	}

	BenchmarkClientGOB(conn)
	BenchmarkClientBinary(conn)
	BenchmarkClientConn(conn)
}

func BenchmarkClientGOB(conn net.Conn) {

	enc := gob.NewEncoder(conn)

	bytes := make([]byte, MessageSize)

	var StartTime, EndTime time.Time
	StartTime = time.Now()
	for i := 0; i < NbMessages; i++ {

		err := enc.Encode(bytes)
		if err != nil {
			log.Panic("error sending message")
		}
	}
	EndTime = time.Now()
	throughput := CalculateThroughput(StartTime, EndTime, MessageSize, NbMessages)
	fmt.Println("Throughput using GOB Package:\t\t" + throughput)
}

func BenchmarkClientBinary(conn net.Conn) {

	buffer := make([]byte, MessageSize)

	var StartTime, EndTime time.Time
	StartTime = time.Now()
	for i := 0; i < NbMessages; i++ {

		err := binary.Write(conn, binary.LittleEndian, buffer)
		if err != nil {
			log.Panic("error sending message")
		}
	}
	EndTime = time.Now()
	throughput := CalculateThroughput(StartTime, EndTime, MessageSize, NbMessages)
	fmt.Println("Throughput using Binary Package:\t" + throughput)
}

func BenchmarkClientConn(conn net.Conn) {

	var StartTime, EndTime time.Time
	StartTime = time.Now()
	buffer := make(Message, MessageSize)
	for i := 0; i < NbMessages; i++ {

		SendConn(conn, buffer)
	}
	EndTime = time.Now()
	throughput := CalculateThroughput(StartTime, EndTime, MessageSize, NbMessages)
	fmt.Println("Throughput using net Package:\t\t" + throughput)
}

func SendConn(conn net.Conn, bytes []byte) {

	ping := 0
	for ping < int(MessageSize) {

		sentBytes, err := conn.Write(bytes[ping:MessageSize])
		if err != nil || sentBytes <= 0 {
			log.Panic("error sending message")
		}
		ping += sentBytes
	}
}
