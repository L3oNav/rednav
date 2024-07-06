package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"rednav/app"
	server "rednav/server"
	"strings"
	"syscall"
)

func main() {
	var replica_of string
	port := flag.Int("port", 3312, "Port to listen on")
	host := flag.String("host", "localhost", "Host to listen on")
	flag.StringVar(&replica_of, "replica_of", "", "Host to replicate from")
	flag.Parse()

	var replicaHost string
	var replicaPort int
	if replica_of != "" {
		replicaParts := strings.Split(replica_of, " ")
		if len(replicaParts) == 2 {
			replicaHost = replicaParts[0]
			fmt.Sscanf(replicaParts[1], "%d", &replicaPort)
		} else {
			fmt.Println("Invalid format for --replicaof. Expected format: host port")
			return
		}
	} else {
		replicaHost = ""
		replicaPort = 0
	}

	config := app.NewConfig(*host, *port, replicaHost, replicaPort)
	vault := app.NewVault(config)

	local_server := server.NewServer(vault, fmt.Sprintf("%s:%d", *host, *port))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		fmt.Printf("Received signal: %s. Shutting down...\n", sig)
		local_server.Shutdown()
	}()

	fmt.Printf("Server running on %s:%d\n", *host, *port)
	if !vault.IsMaster() {
		fmt.Printf("Replicating from %s:%d\n", replicaHost, replicaPort)
		go local_server.HeyListenMaster()
	}

	local_server.HeyListen()
}
