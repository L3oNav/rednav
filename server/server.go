package server

import (
	"fmt"
	"net"
	"os"
	"rednav/app"
	"rednav/commands"
	"rednav/utils"
	"strings"
	"time"
)

var (
	kMaxArgs = 128
	kMaxMsg  = 4096
)

type Server struct {
	listener           net.Listener
	conn               net.Conn
	vault              *app.Vault
	talkingWithMain    bool
	talkingWithReplica bool
	address            string
	masterConn         net.Conn
	Replicas           []*net.Conn
	ReplicasStatus     map[string]bool
	role               string
	quitch             chan struct{}
}

func NewServer(vault *app.Vault, local_addr string) *Server {
	role := "main"
	if !vault.IsMaster() {
		role = "replica"
	}
	server := &Server{
		address:            local_addr,
		Replicas:           make([]*net.Conn, 0),
		ReplicasStatus:     make(map[string]bool),
		talkingWithMain:    false,
		talkingWithReplica: false,
		vault:              vault,
		quitch:             make(chan struct{}),
		role:               role,
	}
	if !server.vault.IsMaster() {
		server.masterConn = server.vault.MasterConn
	}
	return server
}

func (s *Server) ReplicasConnection(address string) {
	s.talkingWithReplica = true
	s.vault.ReplicaPresent = true
	var conn net.Conn
	var err error
	for attempt := 0; attempt < 5; attempt++ {
		conn, err = net.Dial("tcp", address)
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d: Error connecting to replica: %v\n", attempt+1, err)
		time.Sleep(250 * time.Millisecond)
	}
	if err != nil {
		fmt.Println("Error connecting to replica: ", err)
		return
	}
	s.Replicas = append(s.Replicas, &conn)
}

func (s *Server) HeyListenMaster() {
	fmt.Printf("INFO || Listening on %s\n", s.address)
	s.acceptMasterLoop(s.masterConn)
}

func (s *Server) acceptMasterLoop(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, kMaxMsg)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error accepting connection from master: ", err)
			return
		}
		s.handleMasterConnection(buffer[:n])
	}
}

func (s *Server) handleMasterConnection(data []byte) {
	commands, err := utils.ParseReq(data, uint32(len(data)))
	if err != nil {
		fmt.Println("Failed to parse commands from master:", err)
		return
	}
	s.handleMasterCommands(commands)
}

func (s *Server) handleMasterCommands(command []string) {
	cmdName := command[0]
	args := make([]commands.Command, len(command)-1)
	for i, arg := range command[1:] {
		args[i] = commands.Command{Typ: "bulk", Bulk: arg}
	}

	// Process command without sending response
	if handler, exists := commands.Handlers[cmdName]; exists {
		handler(s.vault, args, nil) // nil as ServerActions since we don't need any server actions in this context
	} else {
		fmt.Printf("Unknown command from master: %s\n", cmdName)
	}

}

func (s *Server) HeyListen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		panic(err)
	}
	s.listener = listener
	go s.acceptLoop()
	<-s.quitch
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		s.conn = conn
		fmt.Printf("INFO || Connection accepted, conn %v\n", conn.LocalAddr().String())
		if err != nil {
			os.Exit(1)
		}
		go s.handleConnection(conn)
	}
}

func (sm *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, kMaxMsg)
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		request_parsed, err := utils.ParseReq(buf[:n], uint32(kMaxMsg))
		if err != nil {
			//send the err back to the client
			conn.Write([]byte(err.Error()))
			return
		}

		response := sm.handleCommand(request_parsed)
		conn.Write(response)
		// disconect the client
	}
}

func (s *Server) handleCommand(message []string) []byte {
	result := []byte("-ERR Unknown command\r\n")
	if len(message) == 0 {
		return []byte("-ERR Empty command\r\n")
	}

	cmdName := strings.ToUpper(message[0])
	args := make([]commands.Command, len(message)-1)
	for i, arg := range message[1:] {
		args[i] = commands.Command{Typ: "bulk", Bulk: arg}
	}

	// Check if the command exists in the handlers map
	if handler, exists := commands.Handlers[cmdName]; exists {
		response := handler(s.vault, args, s)
		result = formatResponse(response)
	}

	if s.vault.IsMaster() && isWriteCommand(cmdName) {
		fmt.Printf("INFO || Sending to replicas: %s\n", cmdName)
		s.sendToReplicas(cmdName, args)
	}

	if !s.vault.IsMaster() && isWriteCommand(cmdName) {
		fmt.Printf("INFO || Sending to master: %s\n", cmdName)
		s.sendToMaster(cmdName, args)
	}

	fmt.Printf("INFO || Command Result %s\n", result)
	return result
}

func (s *Server) sendToMaster(cmd string, args []commands.Command) {
	encoded := EncodeCommand(cmd, args)
	response, err := s.vault.MasterConn.Write(encoded)

	if err != nil {
		fmt.Println("Error sending to master: ", err)
	}

	fmt.Printf("Response from master: %d\n", response)
}

// Propagate commands to replicas
func (s *Server) sendToReplicas(cmd string, args []commands.Command) {
	encoded := EncodeCommand(cmd, args)
	for _, replica := range s.Replicas {
		(*replica).SetWriteDeadline(time.Now().Add(5 * time.Second)) // Set deadline to avoid blocking
		_, err := (*replica).Write(encoded)
		if err != nil {
			fmt.Println("Error writing to replica: ", err)
		}
	}
}

func (s *Server) Shutdown() {
	close(s.quitch)
	if s.listener != nil {
		s.listener.Close()
	}
}

func isWriteCommand(cmd string) bool {
	writeCommands := []string{"SET", "DEL"}
	for _, wc := range writeCommands {
		if wc == cmd {
			return true
		}
	}
	return false
}

func EncodeCommand(cmd string, args []commands.Command) []byte {
	var sb strings.Builder

	// Write array header (cmdName + args)
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(args)+1))
	sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(cmd), cmd))

	// Write each argument in bulk string format
	for _, cmd := range args {
		sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(cmd.Bulk), cmd.Bulk))
	}

	return []byte(sb.String())
}

func formatResponse(cmd commands.Command) []byte {
	switch cmd.Typ {
	case "string":
		return []byte(cmd.Str + "\r\n")
	case "list":
		return []byte("*" + fmt.Sprint(len(cmd.List)) + "\r\n" + strings.Join(cmd.List, "\r\n") + "\r\n")
	case "nil":
		return []byte("$-1\r\n")
	case "err":
		return []byte("-ERR " + cmd.Err + "\r\n")
	case "arr", "multi":
		var response []byte
		for _, subCmd := range cmd.Arr {
			response = append(response, formatResponse(subCmd)...)
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(cmd.Arr), response))
	case "bulk":
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(cmd.Bulk), cmd.Bulk))
	default:
		fmt.Printf("Unknown response type: %v\n", cmd)
		return []byte("-ERR Unknown response type\r\n")
	}
}
