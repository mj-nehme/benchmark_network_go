package main

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

var server_conn *net.TCPConn
var bytes []byte
var enc *gob.Encoder

func run_client() {

	completeAddress := IP + ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", completeAddress)
	if err != nil {
		log.Panicf("Error resolving address %s", err)
	}

	new_connection, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Panicf("DialTCP Error %s", err)
	}
	server_conn = new_connection
	enc = gob.NewEncoder(server_conn)

	var fns []func() error
	fns = append(fns, BenchmarkClientGOB)
	fns = append(fns, BenchmarkClientBinary)
	fns = append(fns, BenchmarkClientConn)
	fmt.Println("Throughput: ")

	for _, fn := range fns {

		var StartTime, EndTime time.Time
		StartTime = time.Now()

		bytes = make([]byte, MessageSize)

		for i := 0; i < NbMessages; i++ {
			fn()
			if err != nil {
				log.Panic("error sending message")
			}
		}
		bytes = nil
		EndTime = time.Now()
		throughput := CalculateThroughput(StartTime, EndTime, MessageSize, NbMessages)
		fmt.Println(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name() + ":\t" + throughput)
	}
}

func BenchmarkClientGOB() error {

	return enc.Encode(bytes)
}

func BenchmarkClientBinary() error {

	return binary.Write(server_conn, binary.LittleEndian, bytes)

}

func BenchmarkClientConn() error {

	return SendConn(bytes)
}

func SendConn(bytes []byte) error {

	ping := 0
	for ping < int(MessageSize) {

		sentBytes, err := server_conn.Write(bytes[ping:MessageSize])
		if err != nil || sentBytes <= 0 {
			return fmt.Errorf("error sending message")
		}
		ping += sentBytes
	}
	return nil
}
