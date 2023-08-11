package db

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp/v3"
	"github.com/pkg/errors"
)

func Get(
	host string,
	port int,
	dbname string,
	initCallback func(dbmap *gorp.DbMap),
) (*gorp.DbMap, error) {

	url := fmt.Sprintf("root:@tcp(%s:%d)/%s?parseTime=true", host, port, dbname)

	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "open db")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping")
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{
		Engine:   "InnoDB",
		Encoding: "utf8",
	}}

	initCallback(dbmap)

	if err := dbmap.CreateTablesIfNotExists(); err != nil {
		return nil, errors.Wrap(err, "create tables")
	}

	return dbmap, nil
}
