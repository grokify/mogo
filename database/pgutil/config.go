package pgutil

import (
	"encoding/json"
	"strings"

	"github.com/go-pg/pg"
)

type PgConn struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func NewPgConnJSON(data []byte) (PgConn, error) {
	conn := PgConn{}
	err := json.Unmarshal(data, &conn)
	return conn, err
}

func (conn *PgConn) Address() string {
	if len(conn.Host) > 0 && len(conn.Port) > 0 {
		return "addr=" + conn.Host + ":" + conn.Port
	}
	return ""
}

func (conn *PgConn) Encode() string {
	conn.Trim()
	parts := []string{}
	if len(conn.User) > 0 {
		parts = append(parts, "user="+conn.User)
	}
	if len(conn.Password) > 0 {
		parts = append(parts, "password="+conn.Password)
	}
	if len(conn.Host) > 0 && len(conn.Port) > 0 {
		parts = append(parts, conn.Address())
	}
	if len(conn.Database) > 0 {
		parts = append(parts, "dbname="+conn.Database)
	}
	if len(conn.SSLMode) > 0 {
		parts = append(parts, "sslmode="+conn.SSLMode)
	}
	return strings.Join(parts, " ")
}

func (conn *PgConn) GoPgOptions() *pg.Options {
	return &pg.Options{
		Addr:     conn.Address(),
		User:     conn.User,
		Password: conn.Password,
		Database: conn.Database}
}

func (conn *PgConn) Trim() {
	conn.Host = strings.TrimSpace(conn.Host)
	conn.Port = strings.TrimSpace(conn.Port)
	conn.User = strings.TrimSpace(conn.User)
	conn.Password = strings.TrimSpace(conn.Password)
	conn.Database = strings.TrimSpace(conn.Database)
	conn.SSLMode = strings.TrimSpace(conn.SSLMode)
}
