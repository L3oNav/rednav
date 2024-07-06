package commands

import (
	"fmt"
	"rednav/app"
	"rednav/interfaces"
)

func PSync(vault *app.Vault, cmd []Command, actions interfaces.ServerActions) Command {
	// Add replica connectioon to the server
	resp := make([]Command, 2)
	resp[0] = Command{Typ: "string", Str: fmt.Sprintf("+FULLRESYNC %s %d\r\n", vault.MainReplicaID, vault.MainReplicaOffset)}
	file := vault.RDBParsed()
	resp[1] = Command{Typ: "bulk", Bulk: string(file)}

	return Command{Typ: "multi", Str: "OK", Arr: resp}
}
