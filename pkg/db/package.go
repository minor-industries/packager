package db

import (
	"database/sql"
	"github.com/go-gorp/gorp/v3"
	_ "github.com/go-sql-driver/mysql"
)

type Package struct {
	ID        string
	Name      string
	Major     int
	Minor     int
	Patch     int
	Arch      string
	OS        string
	GitRef    sql.NullString
	Filename  sql.NullString
	Hash      sql.NullString
	Signature sql.NullString
}

func DBMapInit(dbmap *gorp.DbMap) {
	dbmap.AddTable(Package{}).SetKeys(false, "ID")
}
