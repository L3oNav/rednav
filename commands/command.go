package commands

import (
	"rednav/app"
	"rednav/interfaces"
)

type Command struct {
	Typ  string
	Str  string
	Bulk string
	List []string
	Err  string
	Arr  []Command
}

var Handlers = map[string]func(*app.Vault, []Command, interfaces.ServerActions) Command{
	"PING":     Ping,
	"ECHO":     Echo,
	"SET":      Set,
	"GET":      Get,
	"INFO":     Info,
	"REPLCONF": ReplConf,
	"PSYNC":    PSync,
}
