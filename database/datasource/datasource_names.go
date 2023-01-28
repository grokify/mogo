package datasource

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

// dsnBigQuery generates a DSN for BigQuery. Use `Hostname` for `ProjectID`` and `Database` for `locatin/dataset`.
// See https://github.com/go-gorm/bigquery and https://github.com/solcates/go-sql-bigquery .
func dsnBigQuery(ds DataSource) (string, error) {
	// "bigquery://projectid/location/dataset"
	// "bigquery://projectid/dataset"
	dsn := strings.Join(
		stringsutil.SliceCondenseSpace([]string{ds.Hostname, ds.Database}, false, false),
		"/")
	if len(dsn) == 0 {
		return "", errors.New("no dsn information provided")
	}
	schemePlus := "bigquery://"
	if !strings.Contains(dsn, schemePlus) {
		return schemePlus + dsn, nil
	} else if strings.Index(dsn, schemePlus) == 0 {
		return dsn, nil
	}
	return "", errors.New("invalid format")
}

func dsnMySQL(ds DataSource) (string, error) {
	// "root:password1@tcp(127.0.0.1:3306)/test"
	dsn := "tcp(" + ds.Host() + ")"
	if len(ds.Database) > 0 {
		dsn += "/" + ds.Database
	}
	if up := ds.UserPassword(); len(up) > 0 {
		dsn = up + "@" + dsn
	}
	return dsn, nil
}

// dsnPostgres produces a URI DSN connection string
func dsnPostgres(ds DataSource) (string, error) {
	// postgres example: postgres://{user}:{password}@{host}:{port}/{database-name}
	schemaPlus := "postgres://"
	dsn := ds.Host() + "/" + ds.Database
	if un := strings.TrimSpace(ds.UserPassword()); len(un) > 0 {
		dsn = schemaPlus + un + "@" + dsn
	} else {
		dsn = schemaPlus + dsn
	}
	sslmode, err := SSLModeParse(ds.SSLMode)
	if err != nil {
		return dsn, err
	}
	if len(sslmode) > 0 {
		dsn += "?sslmode=" + sslmode
	}
	return dsn, nil
}

func dsnSQLite3(ds DataSource) (string, error) {
	return ds.Database, nil
}
