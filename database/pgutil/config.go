package pgutil

import (
	"encoding/json"
	"strings"

	"github.com/go-pg/pg"
)

type PgConn struct {
	Host     string
	Port     string
	Addr     string `josn:"Addr,omitempty"`
	User     string `json:"User,omitempty"`
	Password string `json:"Password,omitempty"`
	Database string `json:"Database,omitempty"`
	SSLMode  string `json:"sslmode,omitempty"`
}

func NewPgConnJSON(data []byte) (PgConn, error) {
	conn := PgConn{}
	err := json.Unmarshal(data, &conn)
	return conn, err
}

func (conn *PgConn) Address() string {
	if len(conn.Addr) > 0 {
		return conn.Addr
	} else if len(conn.Host) > 0 && len(conn.Port) > 0 {
		return conn.Host + ":" + conn.Port
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
		parts = append(parts, "addr="+conn.Address())
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

func (conn *PgConn) Client() *pg.DB {
	return pg.Connect(conn.GoPgOptions())
}

func (conn *PgConn) Trim() {
	conn.Addr = strings.TrimSpace(conn.Addr)
	conn.Host = strings.TrimSpace(conn.Host)
	conn.Port = strings.TrimSpace(conn.Port)
	conn.User = strings.TrimSpace(conn.User)
	conn.Password = strings.TrimSpace(conn.Password)
	conn.Database = strings.TrimSpace(conn.Database)
	conn.SSLMode = strings.TrimSpace(conn.SSLMode)
}
