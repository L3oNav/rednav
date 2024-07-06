package commands

import (
	"fmt"
	"rednav/app"
	"rednav/interfaces"
	"strings"
)

func ReplConf(vault *app.Vault, args []Command, actions interfaces.ServerActions) Command {

	if len(args) < 2 {
		return Command{Typ: "err", Err: "Not enough arguments for REPLCONF"}
	}

	option := strings.ToUpper(args[0].Bulk)
	switch option {
	case "LISTENING-PORT":
		actions.ReplicasConnection(args[1].Bulk)
		// Handle setting listening port
		if vault.IsMaster() {
			return Command{Typ: "string", Str: "+OK"}
		} else {
			return Command{Typ: "list", List: []string{app.RELPCONF, app.ACK, fmt.Sprint(vault.MainReplicaOffset)}}
		}

	case "CAPA":
		// Handle capabilities
		// Example: Do something with value (e.g., vault.SetCapabilities(value))
		return Command{Typ: "string", Str: "+OK"}
	default:
		return Command{Typ: "err", Err: "Unknown REPLCONF option"}
	}
}
