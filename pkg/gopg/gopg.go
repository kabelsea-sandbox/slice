package gopg

import (
	"net"

	"github.com/go-pg/pg/v9"
)

// NewConnection creates new postgres connection with provided params.
func NewConnection(host, port, user, database, password string) (*pg.DB, error) {
	conn := pg.Connect(&pg.Options{
		Addr:     net.JoinHostPort(host, port),
		User:     user,
		Database: database,
		Password: password,
	})
	_, err := conn.Exec("select 1")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
