package dbutil

import (
	"database/sql"
	"fmt"

	"github.com/dotenx/dotenx/ao-api/config"
	dbpkg "github.com/dotenx/dotenx/ao-api/db"
	"github.com/dotenx/dotenx/ao-api/pkg/utils"
	"github.com/jmoiron/sqlx"
)

/*
As we create a new database for each project, we need a new connection for each project.
To avoid connection leak, we need to close the connection after each query.
As a better solution, we can use a connection pool in the future.

The consumer of GetDbInstance function MUST call fn(db.Connection) to close the connection.

Usage:

func usage_example() {
	db, fn, _ := GetDbInstance("", "")
	defer fn(db.Connection)
	if err := db.Connection.Ping(); err != nil {
		panic(err)
	}
	// execute query
}

*/

type PostQueryCallback func(*sqlx.DB) error

func GetDbInstance(accountId string, projectName string) (*dbpkg.DB, PostQueryCallback, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Configs.Database.User, config.Configs.Database.Password, config.Configs.Database.Host, config.Configs.Database.Port, utils.GetProjectDatabaseName(accountId, projectName))
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, nil, err
	}

	return &dbpkg.DB{
		Connection: sqlx.NewDb(db, "postgres"),
		Driver:     dbpkg.Postgres,
	}, closeAfterQuery, nil

}

func closeAfterQuery(db *sqlx.DB) error {
	return db.Close()
}
