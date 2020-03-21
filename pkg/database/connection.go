// Contains utility functions around the database

package database

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"log"
	"soko/pkg/config"
	"soko/pkg/models"
)

// DBCon is the connection handle
// for the database
var (
	DBCon *pg.DB
)

// CreateSchema creates the tables in the database
// in case they don't alreay exist
func CreateSchema() error {
	for _, model := range []interface{}{(*models.Package)(nil),
		(*models.Category)(nil),
		(*models.Version)(nil),
		(*models.Commit)(nil),
		(*models.KeywordChange)(nil),
		(*models.CommitToPackage)(nil),
		(*models.CommitToVersion)(nil),
		(*models.Useflag)(nil),
		(*models.Application)(nil)} {

		err := DBCon.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}

	}
	return nil
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery is used to log SQL queries
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	log.Println(q.FormattedQuery())
	return nil
}

// Connect is used to connect to the database
// and turn on logging if desired
func Connect() {
	DBCon = pg.Connect(&pg.Options{
		User:     config.PostgresUser(),
		Password: config.PostgresPass(),
		Database: config.PostgresDb(),
		Addr:     config.PostgresHost() + ":" + config.PostgresPort(),
	})

	if config.Debug() == "true" {
		DBCon.AddQueryHook(dbLogger{})
	}

	err := CreateSchema()
	if err != nil {
		panic(err)
	}

}
