package models

import (
	"github.com/jmoiron/sqlx"
)

var (
	Db *sqlx.DB
)

type MysqlConfig struct {
	UserName string
	PassWd string
	Host string
	Port int
	Database string
}

func SetDb(db *sqlx.DB) {
	Db = db
}

