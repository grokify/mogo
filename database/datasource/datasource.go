package datasource

import (
	"database/sql"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/netutil"
)

const (
	DriverBigQuery             = "bigquery"
	DriverGodror               = "godror"
	DriverMySQL                = "mysql"
	DriverOracle               = "oracle"
	DriverPostgres             = "postgres"
	DriverSQLite3              = "sqlite3"
	SchemePostgres             = "postgres"
	DefaultPortMySQL    uint16 = 3306
	DefaultPortPostgres uint16 = 5432
)

type DataSource struct {
	Driver   string              `json:"driver"`
	DSN      string              `json:"dsn"`
	Hostname string              `json:"hostname"` // does not include port
	Port     uint16              `json:"port"`     // 0-65535
	User     string              `json:"user,omitempty"`
	Password string              `json:"password,omitempty"`
	Database string              `json:"database,omitempty"`
	Query    map[string][]string `json:"query,omitempty"`
}

func (ds DataSource) Open() (*sql.DB, error) {
	dsn, err := ds.Name()
	if err != nil {
		return nil, err
	}
	return sql.Open(ds.Driver, dsn)
}

var (
	ErrDriverNameEmpty    = errors.New("driver name cannot be empty")
	ErrDriverNotSupported = errors.New("driver not supported")
)

// Name produces a URI DSN connection string
func (ds *DataSource) Name() (string, error) {
	ds.trim()
	if ds.DSN != "" {
		return ds.DSN, nil
	}
	switch ds.Driver {
	case DriverBigQuery:
		return dsnBigQuery(ds)
	case DriverGodror:
		return dsnGodror(ds)
	case DriverMySQL:
		return dsnMySQL(ds), nil
	case DriverOracle:
		return dsnOracle(ds), nil
	case DriverPostgres:
		return dsnPostgres(ds)
	case DriverSQLite3:
		return dsnSQLite3(ds), nil
	case "":
		return "", ErrDriverNameEmpty
	default:
		return "", errorsutil.Wrapf(ErrDriverNotSupported, "driver name (%s)", ds.Driver)
	}
}

func (ds DataSource) Host() string {
	if port := ds.PortOrDefault(); port > 0 {
		return ds.HostnameOrDefault() + ":" + strconv.Itoa(int(port))
	}
	return strings.TrimSpace(ds.HostnameOrDefault())
}

func (ds DataSource) HostDatabase(sep string) string {
	parts := []string{}
	if host := strings.TrimSpace(ds.Host()); host != "" {
		parts = append(parts, host)
	}
	if db := strings.TrimSpace(ds.Database); db != "" {
		parts = append(parts, db)
	}
	return strings.TrimSpace(strings.Join(parts, sep))
}

// HostDatabaseQuery can be used when there's no need to validate query params
func (ds DataSource) HostDatabaseQuery(sep string) string {
	dsn := ds.HostDatabase(sep)
	if len(ds.Query) > 0 {
		qry := strings.TrimSpace(url.Values(ds.Query).Encode())
		if qry != "" {
			dsn += "?" + qry
		}
	}
	return dsn
}

func (ds DataSource) UserPassHostDatabase() string {
	dsn := ds.HostDatabase("/")
	if un := strings.TrimSpace(ds.UserPassword()); un != "" {
		dsn = un + "@" + dsn
	}
	return dsn
}

// UserPassHostDatabaseQuery can be used when there's no need to validate query params
func (ds DataSource) UserPassHostDatabaseQuery() string {
	dsn := ds.UserPassHostDatabase()
	if len(ds.Query) > 0 {
		qry := strings.TrimSpace(url.Values(ds.Query).Encode())
		if qry != "" {
			dsn += "?" + qry
		}
	}
	return dsn
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
	if strings.TrimSpace(ds.Hostname) == "" {
		return netutil.HostLocalhost
	}
	return strings.TrimSpace(ds.Hostname)
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
	ds.Driver = strings.ToLower(strings.TrimSpace(ds.Driver))
	ds.DSN = strings.TrimSpace(ds.DSN)
	ds.Hostname = strings.TrimSpace(ds.Hostname)
	ds.User = strings.TrimSpace(ds.User)
	ds.Password = strings.TrimSpace(ds.Password)
	ds.Database = strings.TrimSpace(ds.Database)
}
