package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
)

// Helper function to send a command and read the response
func sendCommand(conn net.Conn, command string) (string, error) {
	fmt.Fprintf(conn, command)
	response, err := bufio.NewReader(conn).ReadString('\n')
	return response, err
}

// Test Command Propagation
func TestCommandPropagation(t *testing.T) {
	replicaConn, err := net.Dial("tcp", "localhost:3313")
	if err != nil {
		t.Fatalf("Failed to connect to server as replica: %v", err)
	}
	defer replicaConn.Close()

	clientConn, err := net.Dial("tcp", "localhost:3312")
	if err != nil {
		t.Fatalf("Failed to connect to server as client: %v", err)
	}
	defer clientConn.Close()

	// Send write commands as a client
	commands := []string{
		"*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$1\r\n1\r\n",
		"*3\r\n$3\r\nSET\r\n$3\r\nbar\r\n$1\r\n2\r\n",
		"*3\r\n$3\r\nSET\r\n$3\r\nbaz\r\n$1\r\n3\r\n",
	}

	for _, cmd := range commands {
		_, err := sendCommand(clientConn, cmd)
		if err != nil {
			t.Fatalf("Failed to send command %s: %v", cmd, err)
		}
	}

	commands = []string{
		"*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n",
		"*2\r\n$3\r\nGET\r\n$3\r\nbar\r\n",
		"*2\r\n$3\r\nGET\r\n$3\r\nbaz\r\n",
	}

	// Read from the replica and check if the commands are propagated
	for _, cmd := range commands {
		response, err := sendCommand(clientConn, cmd)
		if err != nil {
			t.Fatalf("Failed to send command %s: %v", cmd, err)
		}
		//INFO || Command: GET, Result: {nil   []  []} is not nil
		if strings.Contains(response, "nil") {
			t.Fatalf("Failed to propagate command %s: %v", cmd, err)
		}
	}
}
