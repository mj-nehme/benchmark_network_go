package main

import (
	"fmt"
	"os"
	"strconv"
)

var (
	IP          string = "127.0.0.1"
	MessageSize int    = 1000000
	NbMessages  int    = 1000
	port        int    = 5555
)

type Message []byte

func main() {

	switch len(os.Args) {
	case 2:

		role := os.Args[1]
		switch role {
		case "server":
			run_server()
		case "client":
			run_client()
		default:
			panic("Unknown argument role")
		}
	case 3:

		new_port, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading port %d\nFatal error: %s\n", port, err.Error())
		}
		port = new_port
		role := os.Args[1]
		switch role {
		case "server":
			run_server()
		case "client":
			run_client()
		default:
			panic("Unknown argument role")
		}
	default:

		panic("Unknown arguments")
	}
}
