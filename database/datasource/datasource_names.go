package datasource

import (
	"errors"
	"net/url"
	"strings"

	"github.com/grokify/mogo/log/logutil"
)

// dsnBigQuery generates a DSN for BigQuery. Use `Hostname` for `ProjectID“ and `Database` for `location/dataset`.
// See https://github.com/go-gorm/bigquery and https://github.com/solcates/go-sql-bigquery .
func dsnBigQuery(ds *DataSource) (string, error) {
	// "bigquery://projectid/location/dataset"
	// "bigquery://projectid/dataset"
	dsn := ds.HostDatabase("/")
	if len(dsn) == 0 {
		return "", errors.New("no dsn information provided")
	}
	return "bigquery://" + dsn, nil
}

// dsnGodror builds a logfmt string.
func dsnGodror(ds *DataSource) (string, error) {
	// https://github.com/godror/godror
	data := map[string][]string{
		"connectString": {ds.HostDatabaseQuery("/")}, // aka Oracle Easy Connect String https://download.oracle.com/ocomdocs/global/Oracle-Net-19c-Easy-Connect-Plus.pdf
	}
	if u := strings.TrimSpace(ds.User); u != "" {
		data["user"] = []string{u}
	}
	if p := strings.TrimSpace(ds.Password); p != "" {
		data["password"] = []string{p}
	}
	return logutil.LogfmtString(data)
}

func ExampleGodrorQueryParams() map[string][]string {
	// see: https://github.com/godror/godror
	return map[string][]string{
		"connect_timeout": {"15"},
	}
}

func dsnMySQL(ds *DataSource) string {
	// "root:password1@tcp(127.0.0.1:3306)/test"
	dsn := "tcp(" + ds.Host() + ")"
	if db := strings.TrimSpace(ds.Database); db != "" {
		dsn += "/" + db
	}
	if up := ds.UserPassword(); up != "" {
		dsn = up + "@" + dsn
	}
	return dsn
}

func dsnOracle(ds *DataSource) string {
	// see: https://docs.oracle.com/database/121/ODPNT/featConnecting.htm
	// `host:port?query` is aka Oracle Easy Connect String https://download.oracle.com/ocomdocs/global/Oracle-Net-19c-Easy-Connect-Plus.pdf
	return "oracle://" + ds.UserPassHostDatabaseQuery()
}

func ExampleOracleQueryParams() map[string][]string {
	// see: https://blogs.oracle.com/developers/post/connecting-a-go-application-to-oracle-database
	// connectionString := "oracle://" + dbParams["username"] + ":" + dbParams["password"] + "@" + dbParams["server"] + ":" + dbParams["port"] + "/" + dbParams["service"]; if val, ok := dbParams["walletLocation"]; ok && val != "" {connectionString += "?TRACE FILE=trace.log&SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(dbParams["walletLocation"])
	return map[string][]string{
		"TRACE FILE": {"trace.log"},
		"SSL":        {"enable"},
		"SSL Verify": {"false"},
		"WALLET":     {"myWalletLocation"},
	}
}

/*
oconnectionStringnnectionString := "oracle://" + dbParams["username"] + ":" + dbParams["password"] + "@" + dbParams["server"] + ":" + dbParams["port"] + "/" + dbParams["service"]

if val, ok := dbParams["walletLocation"]; ok && val != "" {
   connectionString += "?TRACE FILE=trace.log&SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(dbParams["walletLocation"])
}
db, err := sql.Open("oracle", connectionString)
*/

// dsnPostgres produces a URI DSN connection string
func dsnPostgres(ds *DataSource) (string, error) {
	// postgres example: postgres://{user}:{password}@{host}:{port}/{database-name}
	dsn := "postgres://" + ds.UserPassHostDatabase()
	if len(ds.Query) > 0 {
		q := url.Values(ds.Query)
		sslmode := strings.TrimSpace(q.Get(PgSSLModeParam))
		if sslmode != "" {
			sslmodeCan, err := SSLModeParse(sslmode)
			if err != nil {
				return dsn, err
			}
			if sslmodeCan != sslmode {
				q.Set(PgSSLModeParam, sslmodeCan)
			}
		}
		qry := strings.TrimSpace(q.Encode())
		if qry != "" {
			dsn += "?" + qry
		}
	}
	return dsn, nil
}

func dsnSQLite3(ds *DataSource) string {
	return ds.Database
}
