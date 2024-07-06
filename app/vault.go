package app

import (
	"encoding/base64"
	"fmt"
	"net"
	"rednav/utils"
	"sync"
	"time"
)

type Vault struct {
	config                 *Config
	role                   string
	MainReplicaID          string
	MainReplicaOffset      int
	memory                 *MemoryStorage
	ReplicaPresent         bool
	alreadyConnectedMaster bool
	MasterConn             net.Conn
	mutex                  sync.Mutex
}

const (
	ACK            = "ACK"
	CAPABILITY     = "capa"
	CONFIG         = "config"
	ECHO           = "echo"
	GET            = "get"
	TYPE           = "type"
	GETACK         = "GETACK"
	MASTER         = "master"
	INFO           = "info"
	LISTENING_PORT = "listening-port"
	SET            = "set"
	XADD           = "xadd"
	REPLICA        = "replica"
	PING           = "ping"
	PX             = "px"
	PSYNC          = "PSYNC"
	REPLICATION    = "replication"
	RELPCONF       = "REPLCONF"

	OK = "+OK\r\n"

	LEN_CAPABILITY     = 2
	LEN_CONFIG         = 1
	LEN_ECHO           = 2
	LEN_GET            = 2
	LEN_TYPE           = 2
	LEN_INFO           = 2
	LEN_LISTENING_PORT = 2
	LEN_REPL_CONF      = 1
	LEN_SET            = 3
	LEN_PING           = 1
	LEN_PX             = 2
	LEN_PSYNC          = 3
	LEN_GETACK         = 2
	LEN_XADD           = 2
)

func NewVault(c *Config) *Vault {
	v := &Vault{
		memory:                 NewMemoryStorage(),
		ReplicaPresent:         false,
		alreadyConnectedMaster: false,
		config:                 c,
	}
	if c.Master_host == "" && c.Master_port == 0 {
		v.role = MASTER
		v.MainReplicaID = utils.GenerateAlphanumericString()
	} else {
		v.role = REPLICA
		v.MainReplicaID = ""
		v.MainReplicaOffset = 0
	}
	if v.role == REPLICA {
		v.MasterConn = v.OpenConnectionToMaster()
	}
	return v
}

func (v *Vault) SetXAdd(setVals []map[string][]string, data map[string][]interface{}) string {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	var lifetime *time.Time
	for i, item := range setVals {
		for key, value := range item {
			if px, exists := data["px"]; exists {
				ms := int(px[i].(float64))
				t := time.Now().Add(time.Duration(ms) * time.Millisecond)
				lifetime = &t
			}
			v.memory.Save(key, value, lifetime)
		}
	}
	return OK
}

func (v *Vault) SetMemory(key string, value string, expiration *time.Time) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.memory.Save(key, value, expiration)
}

func (v *Vault) GetMemory(key string) interface{} {
	return v.memory.Get(key)
}

func (v *Vault) GetType(key string) string {
	return v.memory.GetType(key)
}

func (v *Vault) GetConfig() *Config {
	return v.config
}

func (v *Vault) GetInfo() string {
	if v.IsMaster() {
		return fmt.Sprintf("role:%s\n", v.role)
	} else {
		return fmt.Sprintf("role:%s\nmain_replid:%s\nmain_repl_offset:%d\n", v.role, v.MainReplicaID, v.MainReplicaOffset)
	}

}

func (v *Vault) OpenConnectionToMaster() net.Conn {
	address := fmt.Sprintf("%s:%d", v.config.Master_host, v.config.Master_port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	// Perform the handshake only once
	v.DoHandshake(conn)
	fmt.Printf("Connected to master %s\n", address)
	return conn
}

func (v *Vault) DoHandshake(conn net.Conn) {
	address := fmt.Sprintf("%s:%d", v.config.Host, v.config.Port)
	// Do the handshake
	conn.Write([]byte("*1\r\n$4\r\nping\r\n"))
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read from master")
	}
	fmt.Printf("Received: %s\n", string(buf[:n]))
	time.Sleep(100 * time.Millisecond)
	conn.Write([]byte(fmt.Sprintf("*4\r\n$8\r\nREPLCONF\r\n$7\r\nLISTENING-PORT\r\n$%d\r\n%s\r\n", len(address), address)))
	buf = make([]byte, 4096)
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read from master")
	}
	fmt.Printf("Received: %s\n", string(buf[:n]))
	time.Sleep(100 * time.Millisecond)
	conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
	buf = make([]byte, 4096)
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read from master")
	}
	fmt.Printf("Received: %s\n", string(buf[:n]))
	time.Sleep(100 * time.Millisecond)
	conn.Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))
	buf = make([]byte, 4096)
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read from master")
	}

	fmt.Printf("Received: %s\n", string(buf[:n]))
	// Return the connection which will be used further
}

func (v *Vault) RDBParsed() string {
	rdb, _ := base64.StdEncoding.DecodeString("UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog==")
	length := len(rdb)
	parsed := fmt.Sprintf("$%d\r\n%s\r\n", length, rdb)
	return parsed
}

func (v *Vault) IsMaster() bool {
	return v.role == MASTER
}
