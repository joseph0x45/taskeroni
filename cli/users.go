package cli

import (
	"fmt"

	"github.com/joseph0x45/taskeroni/db"
	"github.com/joseph0x45/taskeroni/internal/models"
	"github.com/joseph0x45/taskeroni/utils"
	"github.com/teris-io/shortid"
)

func handleUsers(opts *CLIOptions, conn *db.Conn) {
	if opts.Create {
		createUser(opts, conn)
    return
	}
}

func createUser(opts *CLIOptions, conn *db.Conn) {
	if opts.Username == "" {
		fmt.Println("username is required")
		return
	}
	if opts.Password == "" {
		fmt.Println("password is required")
		return
	}
	user, err := conn.GetUserByUsername(opts.Username)
	if err != nil {
		fmt.Println(err)
		return
	}
	if user != nil {
		fmt.Println("Username", opts.Username, "is already taken")
		return
	}
	hash, err := utils.HashPassword(opts.Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	newUser := &models.User{
		ID:       shortid.MustGenerate(),
		Username: opts.Username,
		Password: hash,
	}
	err = conn.InsertUser(newUser)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("User", opts.Username, "created")
}
