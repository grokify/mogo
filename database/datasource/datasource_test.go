package datasource

import (
	"testing"

	"github.com/grokify/mogo/net/netutil"
)

/*
Example without hostname:

"root:root@/my_db?charset=utf8" - https://pkg.go.dev/github.com/astaxie/beego/orm#section-readme
*/

var dsnTests = []struct {
	v       DataSource
	want    string
	wantErr error
}{
	{DataSource{
		Driver:   DriverBigQuery,
		DSN:      "bigquery://projectid/dataset",
		Database: "mydb",
	}, "bigquery://projectid/dataset", nil},
	{DataSource{
		Driver:   DriverBigQuery,
		Hostname: "myprojectid",
		Database: "mydataset",
	}, "bigquery://myprojectid/mydataset", nil},
	{DataSource{
		Driver:   DriverBigQuery,
		Hostname: "myprojectid",
		Database: "mylocation/mydataset",
	}, "bigquery://myprojectid/mylocation/mydataset", nil},
	{DataSource{
		Driver:   DriverMySQL,
		Hostname: netutil.HostLoopbackIPv4,
		Port:     5432,
		User:     "myuser",
		Password: "mypass",
		Database: "mydb",
	}, "myuser:mypass@tcp(127.0.0.1:5432)/mydb", nil},
	{DataSource{
		Driver:   DriverPostgres,
		Hostname: netutil.HostLoopbackIPv4,
		Port:     5432,
		User:     "myuser",
		Password: "mypass",
		Database: "mydb",
	}, "postgres://myuser:mypass@127.0.0.1:5432/mydb", ErrSSLModeNotSUpported},
	{DataSource{
		Driver:   DriverPostgres,
		DSN:      "  postgres://myuser:mypass@127.0.0.1:8888/mydb?sslmode=require",
		Hostname: netutil.HostLoopbackIPv4,
		Port:     5432, // override this since DSN exists
		User:     "myuser",
		Password: "mypass",
		Database: "mydb",
		SSLMode:  SSLModeRequire,
	}, "postgres://myuser:mypass@127.0.0.1:8888/mydb?sslmode=require", nil},
	{DataSource{
		Driver:   DriverPostgres,
		Hostname: netutil.HostLoopbackIPv4,
		Port:     5432,
		User:     "  myuser",
		Password: "mypass",
		Database: "mydb",
		SSLMode:  SSLModeRequire,
	}, "postgres://myuser:mypass@127.0.0.1:5432/mydb?sslmode=require", nil},
	{DataSource{
		Driver:   DriverSQLite3,
		Hostname: netutil.HostLoopbackIPv4,
		Port:     5432,
		User:     "myuser",
		Password: "mypass",
		Database: "./mydb.sql",
	}, "./mydb.sql", nil},
}

func TestDataSourceName(t *testing.T) {
	for _, tt := range dsnTests {
		got, err := tt.v.Name()
		if err != nil {
			if err != tt.wantErr {
				t.Errorf("err DataSource.Name(): err (%s)", err.Error())
			}
		} else {
			if got != tt.want {
				t.Errorf("mismatch DataSource.Name(): want (%s) got (%s)", tt.want, got)
			}
		}
	}
}
