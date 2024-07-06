package commands

import (
	"rednav/app"
	"rednav/interfaces"
)

func Get(v *app.Vault, args []Command, actions interfaces.ServerActions) Command {
	if len(args) < 1 {
		return Command{Typ: "error", Err: "ERR wrong number of arguments for 'get' command"}
	}
	key := args[0].Bulk

	// Retrieve the value from the vault, interfaced return as string.
	value := v.GetMemory(key)

	if value == nil {
		return Command{Typ: "nil"}
	}

	return Command{Typ: "string", Str: value.(string)}
}
