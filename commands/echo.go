package commands

import (
	"rednav/app"
	"rednav/interfaces"
)

func Echo(v *app.Vault, args []Command, actions interfaces.ServerActions) Command {
	if len(args) == 0 {
		return Command{Typ: "string", Str: ""}
	}
	return Command{Typ: "string", Str: args[0].Bulk}
}
