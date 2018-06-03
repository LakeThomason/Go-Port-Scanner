package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var ipAddress string
var portNumberStart, portNumberEnd, timeout int
var wg sync.WaitGroup

func main() {
	// Passed in arguments
	ipAddress = os.Args[1]
	portNumberStart, _ = strconv.Atoi(os.Args[2])
	portNumberEnd, _ = strconv.Atoi(os.Args[3])
	timeout, _ = strconv.Atoi(os.Args[4])

	for i := portNumberStart; i <= portNumberEnd; i++ {
		wg.Add(1)            // Add a process to the wait group
		go dialConnection(i) // Start the goroutine
	}
	wg.Wait() // Wait for all routines to finish
}

func dialConnection(i int) {
	fullAddress := ipAddress + ":" + strconv.Itoa(i)
	// Dial the network:port we want to examine
	conn, err := net.DialTimeout("tcp", fullAddress, time.Duration(timeout)*time.Millisecond)
	if err != nil { // If the port is not open or is filtered
		fmt.Println("Dial error ", err.Error())
	} else { // The port is open
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n") // Make the request by writing to the conn
		var buf bytes.Buffer
		io.Copy(&buf, conn) // Copy the whole contents of the response
		fmt.Println(fullAddress + " responded with " + strconv.Itoa(buf.Len()) + " bytes")
		conn.Close() // Close the connection
	}
	wg.Done()
}
