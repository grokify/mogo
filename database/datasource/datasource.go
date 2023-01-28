package datasource

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/grokify/mogo/net/netutil"
)

const (
	DriverBigQuery             = "bigquery"
	DriverMySQL                = "mysql"
	DriverPostgres             = "postgres"
	DriverSQLite3              = "sqlite3"
	SchemePostgres             = "postgres"
	DefaultPortMySQL    uint16 = 3306
	DefaultPortPostgres uint16 = 5432
)

type DataSource struct {
	Driver   string `json:"driver"`
	DSN      string `json:"dsn"`
	Hostname string `json:"hostname"` // does not include port
	Port     uint16 `json:"port"`     // 0-65535
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	SSLMode  string `json:"sslmode,omitempty"`
}

func (ds DataSource) Open() (*sql.DB, error) {
	dsn, err := ds.Name()
	if err != nil {
		return nil, err
	}
	return sql.Open(ds.Driver, dsn)
}

// Name produces a URI DSN connection string
func (ds DataSource) Name() (string, error) {
	ds.trim()
	if ds.DSN != "" {
		return ds.DSN, nil
	}
	driverName := strings.ToLower(strings.TrimSpace(ds.Driver))
	switch driverName {
	case DriverBigQuery:
		return dsnBigQuery(ds)
	case DriverMySQL:
		return dsnMySQL(ds)
	case DriverPostgres:
		return dsnPostgres(ds)
	case DriverSQLite3:
		return dsnSQLite3(ds)
	default:
		return "", errors.New("db driver not supported")
	}
}

func (ds DataSource) Host() string {
	if port := ds.PortOrDefault(); port > 0 {
		return ds.HostnameOrDefault() + ":" + strconv.Itoa(int(port))
	}
	return ds.HostnameOrDefault()
}

func (ds DataSource) UserPassword() string {
	if ds.User == "" && ds.Password == "" {
		return ""
	} else if ds.Password == "" {
		return ds.User
	}
	return ds.User + ":" + ds.Password
}

func (ds DataSource) HostnameOrDefault() string {
	if len(ds.Hostname) == 0 {
		return netutil.HostLocalhost
	}
	return ds.Hostname
}

func (ds DataSource) PortOrDefault() uint16 {
	if ds.Port == 0 {
		switch ds.Driver {
		case DriverMySQL:
			return DefaultPortMySQL
		case DriverPostgres:
			return DefaultPortPostgres
		}
	}
	return ds.Port
}

func (ds *DataSource) trim() {
	ds.DSN = strings.TrimSpace(ds.DSN)
	ds.Hostname = strings.TrimSpace(ds.Hostname)
	ds.User = strings.TrimSpace(ds.User)
	ds.Password = strings.TrimSpace(ds.Password)
	ds.Database = strings.TrimSpace(ds.Database)
	ds.SSLMode = strings.TrimSpace(ds.SSLMode)
}
