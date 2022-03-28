package datasource

import (
	"fmt"
	"strconv"
	"strings"
)

type DataSource struct {
	Scheme   string
	Host     string
	Port     uint
	Addr     string `json:"Addr,omitempty"`
	User     string `json:"User,omitempty"`
	Password string `json:"Password,omitempty"`
	Database string `json:"Database,omitempty"`
	SSLMode  string `json:"sslmode,omitempty"`
}

func (ds *DataSource) DSN() (string, error) {
	conn := fmt.Sprintf("%s://%s@%s/%s",
		ds.Scheme,
		ds.UserPassword(),
		ds.Hostname(),
		ds.Database)
	sslmode, err := SSLModeParse(ds.SSLMode)
	if err != nil {
		return conn, err
	}
	if len(sslmode) > 0 {
		conn += "?sslmode=" + sslmode
	}
	return conn, nil
}

func (ds *DataSource) Hostname() string {
	return ds.HostOrDefault() + ":" + strconv.Itoa(ds.PortOrDefault())
}

func (ds *DataSource) UserPassword() string {
	if len(ds.User) == 0 && len(ds.Password) == 0 {
		return ""
	}
	return ds.User + ":" + ds.Password
}

func (ds *DataSource) HostOrDefault() string {
	host := ds.Host
	if len(host) == 0 {
		return "localhost"
	}
	return host
}

func (ds *DataSource) PortOrDefault() int {
	if ds.Port != 0 {
		return int(ds.Port)
	}
	return 5432
}

func (ds *DataSource) Trim() {
	ds.Addr = strings.TrimSpace(ds.Addr)
	ds.Host = strings.TrimSpace(ds.Host)
	ds.User = strings.TrimSpace(ds.User)
	ds.Password = strings.TrimSpace(ds.Password)
	ds.Database = strings.TrimSpace(ds.Database)
	ds.SSLMode = strings.TrimSpace(ds.SSLMode)
}
