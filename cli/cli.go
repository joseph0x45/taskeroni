package cli

import "github.com/joseph0x45/taskeroni/db"

type CLIOptions struct {
	Users    bool
	Create   bool
	Username string
	Password string
}

func DispatchCommands(opts *CLIOptions, conn *db.Conn) {
	if opts.Users {
		handleUsers(opts, conn)
		return
	}
}
