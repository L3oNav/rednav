package commands

import (
	"net"
	"testing"
)

func TestPingIntegration(t *testing.T) {
	// Start your server here if needed
	// Create a TCP connection to the server
	conn, err := net.Dial("tcp", "localhost:3312")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send the PING command
	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		t.Fatalf("Failed to write to connection: %v", err)
	}

	// Read the response from the server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	// Check the response
	response := string(buffer[:n])
	expectedResponse := "+PONG"
	if response != expectedResponse {
		t.Errorf("Unexpected response. Got: %s, want: %s", response, expectedResponse)
	}
}
