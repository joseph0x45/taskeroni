package db

import "github.com/joseph0x45/sad"

var migrations = []sad.Migration{
	{
		Version: 1,
		Name:    "create_db_schema",
		SQL: `
      create table users (
        id text not null primary key,
        username text not null unique,
        password text not null
      );

      create table tasks (
        id text not null primary key,
        label text not null,
        status text not null,
        due_date text not null,
        owner_id text not null references users(id)
      );
    `,
	},
}
