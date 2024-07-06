package main

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestConcurrentPings(t *testing.T) {
	const (
		numClients = 100
		address    = "localhost:3312"
		pingCmd    = "*1\r\n$4\r\nPING\r\n"
	)

	var wg sync.WaitGroup
	wg.Add(numClients)

	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", address)
			if err != nil {
				t.Errorf("Client %d: Error connecting: %v", clientID, err)
				return
			}
			defer conn.Close()

			_, err = conn.Write([]byte(pingCmd))
			if err != nil {
				t.Errorf("Client %d: Error writing to connection: %v", clientID, err)
				return
			}

			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				t.Errorf("Client %d: Error reading from connection: %v", clientID, err)
			} else {
				response := string(buffer[:n])
				fmt.Printf("Client %d received response: %s\n", clientID, response)
			}
		}(i)
	}

	wg.Wait()
}
