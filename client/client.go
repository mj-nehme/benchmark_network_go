package main

import (
	"Baseline/Benchmark"
	"Harmony/Errors"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Message []byte

var ip string = "127.0.0.1"
var MessageSize int = 1000000
var NbMessages int = 1000
var port int = 5555

func main() {

	if len(os.Args) > 2 {

		port, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading port %d\nFatal error: %s\n", port, err.Error())
		}
	}
	completeAddress := ip + ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", completeAddress)
	if err != nil {
		Errors.Error(err, "Error resolving address")
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Errors.Error(err, "DialTCP Error")
	}

	BenchmarkGOB(conn)
	BenchmarkBinary(conn)
	BenchmarkConn(conn)
}

func BenchmarkGOB(conn net.Conn) {

	enc := gob.NewEncoder(conn)

	bytes := make([]byte, MessageSize)

	var StartTime, EndTime time.Time
	StartTime = time.Now()
	for i := 0; i < NbMessages; i++ {

		err := enc.Encode(bytes)
		if err != nil {
			panic("error sending message")
		}
	}
	EndTime = time.Now()
	throughput := Benchmark.CalculateThroughput(StartTime, EndTime, MessageSize, NbMessages)
	fmt.Println("Throughput using GOB Package:\t\t" + throughput)
}

func BenchmarkBinary(conn net.Conn) {

	buffer := make([]byte, MessageSize)

	var StartTime, EndTime time.Time
	StartTime = time.Now()
	for i := 0; i < NbMessages; i++ {

		err := binary.Write(conn, binary.LittleEndian, buffer)
		if err != nil {
			panic("error sending message")
		}
	}
	EndTime = time.Now()
	throughput := Benchmark.CalculateThroughput(StartTime, EndTime, MessageSize, NbMessages)
	fmt.Println("Throughput using Binary Package:\t" + throughput)
}

func BenchmarkConn(conn net.Conn) {

	var StartTime, EndTime time.Time
	StartTime = time.Now()
	buffer := make(Message, MessageSize)
	for i := 0; i < NbMessages; i++ {

		SendConn(conn, buffer)
	}
	EndTime = time.Now()
	throughput := Benchmark.CalculateThroughput(StartTime, EndTime, MessageSize, NbMessages)
	fmt.Println("Throughput using net Package:\t\t" + throughput)
}

func SendConn(conn net.Conn, bytes []byte) {

	ping := 0
	for ping < int(MessageSize) {

		sentBytes, err := conn.Write(bytes[ping:MessageSize])
		if err != nil || sentBytes <= 0 {
			panic("error sending message")
		}
		ping += sentBytes
	}
}
