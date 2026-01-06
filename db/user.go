package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/joseph0x45/taskeroni/internal/models"
)

func (c *Conn) getUser(by, value string) (*models.User, error) {
	const queryByID = "select * from users where id=?"
	const queryByUsername = "select * from users where username=?"
	var err error
	user := &models.User{}
	if by == "id" {
		err = c.db.Get(user, queryByID, value)
	} else {
		err = c.db.Get(user, queryByUsername, value)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting user by %s: %w", by, err)
	}
	return user, nil
}

func (c *Conn) GetUserByUsername(username string) (*models.User, error) {
	return c.getUser("username", username)
}

func (c *Conn) GetUserByID(id string) (*models.User, error) {
	return c.getUser("id", id)
}

func (c *Conn) InsertUser(user *models.User) error {
	const query = `
    insert into users (
      id, username, password
    )
    values (
      :id, :username, :password
    );
  `
	_, err := c.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}

func (c *Conn) ChangeUserPassword(id, password string) error {
	const query = "update users set password=? where id=?"
	_, err := c.db.Exec(query, password, id)
	if err != nil {
		return fmt.Errorf("Error while updating user password: %w", err)
	}
	return nil
}
