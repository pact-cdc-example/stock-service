package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type NewPostgresOpts struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func New(opts *NewPostgresOpts) *sql.DB {
	db, err := sql.Open("postgres", createDSNFromOpts(opts))
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func createDSNFromOpts(opts *NewPostgresOpts) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		opts.Host, opts.Port, opts.Username, opts.Password, opts.DBName)
}
