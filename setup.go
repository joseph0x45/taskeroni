package main

import (
	"os/user"

	"github.com/joseph0x45/taskeroni/utils"
)

func setup() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	if err := utils.EnsureDirExists(utils.AppDataDir(user), 0755); err != nil {
		panic(err)
	}
	if err := utils.EnsureDirExists(utils.AppConfigDir(user), 0755); err != nil {
		panic(err)
	}
	return utils.AppDatabasePath(user)
}
