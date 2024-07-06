package commands

import (
	"fmt"
	"rednav/app"
	"rednav/interfaces"
	"strconv"
	"time"
)

func Set(v *app.Vault, args []Command, actions interfaces.ServerActions) Command {
	if len(args) < 2 {
		return Command{Typ: "error", Err: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].Bulk
	value := args[1].Bulk

	var expiration *time.Time

	// Check for optional "PX" argument for expiration time

	if len(args) >= 4 && args[2].Bulk == "px" {
		ms, err := strconv.Atoi(args[3].Bulk)
		if err != nil {
			return Command{Typ: "error", Err: "ERR invalid expiration time"}
		}
		exp := time.Now().Add(time.Duration(ms) * time.Millisecond)
		expiration = &exp
	}

	fmt.Printf("Cmd: SET key=%s, value=%s, expiration=%v\n", key, value, expiration)
	// Store the key-value pair using SetMemory
	v.SetMemory(key, value, expiration)

	return Command{Typ: "string", Str: "+OK"}
}
