package cli

import "github.com/joseph0x45/taskeroni/db"

type CLIOptions struct {
	Username string
	Password string
}

func DispatchCommands(opts *CLIOptions, conn *db.Conn) {
}
