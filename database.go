package main

import (
	"context"
	"database/sql"
	"errors"
)

var ddls = []string{
	`CREATE TABLE IF NOT EXISTS whitelist (
       pubkey text NOT NULL,
       level integer NOT NULL);`,
}

func getWhitelistLevel(ctx context.Context, pubkey string) (int, error) {
	var level int
	err := db.DB.GetContext(ctx, &level, `
        SELECT level FROM whitelist WHERE pubkey = ?
    `, pubkey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return level, nil
}
