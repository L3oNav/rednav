package commands

import (
	"rednav/app"
	"rednav/interfaces"
)

func Info(v *app.Vault, args []Command, actions interfaces.ServerActions) Command {
	//return role of the server
	return Command{Typ: "string", Str: v.GetInfo()}
}
