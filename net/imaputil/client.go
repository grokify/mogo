package imaputil

import (
	"os"
	"strconv"
	"strings"

	"github.com/emersion/go-imap/client"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	GmailIMAPHostname string = "imap.gmail.com"
	GmailIMAPPort     int    = 993
	GmailSMTPHostname string = "smtp.gmail.com"
	DefaultEnvPrefix  string = "IMAP_"
)

type ClientMore struct {
	Config ClientConfig
	Client *client.Client
}

func NewClientMoreEnv(prefix string) (*ClientMore, error) {
	cm := &ClientMore{}
	cc, err := NewClientConfigEnv(prefix)
	if err != nil {
		return nil, err
	}
	cm.Config = cc
	return cm, nil
}

type ClientConfig struct {
	Hostname    string
	Username    string
	Password    string
	Port        int
	TLSRequired bool
}

func NewClientConfigEnv(prefix string) (ClientConfig, error) {
	cc := ClientConfig{
		Hostname:    os.Getenv(prefix + "HOSTNAME"),
		Username:    os.Getenv(prefix + "USERNAME"),
		Password:    os.Getenv(prefix + "PASSWORD"),
		TLSRequired: stringsutil.ToBool(os.Getenv(prefix + "TLS_REQUIRED"))}
	portStr := os.Getenv(prefix + "PORT")
	if len(portStr) > 0 {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return cc, err
		}
		cc.Port = port
	}
	return cc, nil
}

func (cm *ClientMore) ConnectAndLogin() error {
	err := cm.Connect()
	if err != nil {
		return err
	}
	return cm.Login()
}

func (cm *ClientMore) Connect() error {
	clt, err := client.DialTLS(Host(cm.Config.Hostname, cm.Config.Port), nil)
	if err != nil {
		return err
	}
	cm.Client = clt
	return nil
}

func (cm *ClientMore) Login() error {
	return cm.Client.Login(
		cm.Config.Username,
		cm.Config.Password)
}

func (cm *ClientMore) Logout() {
	if cm.Client != nil {
		cm.Client.Logout()
	}
}

func Host(hostname string, port int) string {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) > 0 && port > 0 {
		return hostname + ":" + strconv.Itoa(port)
	} else if len(hostname) > 0 {
		return hostname
	}
	return ""
}

/*
// https://www.lifewire.com/what-are-the-gmail-imap-settings-1170852

Gmail IMAP server address: imap.gmail.com
Gmail IMAP username: Your full Gmail address (for example, example@gmail.com)
Gmail IMAP password: Your Gmail password (use an application-specific Gmail password if you enabled 2-step authentication for Gmail)
Gmail IMAP port: 993
Gmail IMAP TLS/SSL required: yes


https://stackoverflow.com/questions/49009432/emersion-go-imap-how-to-retrieve-and-list-unseen-messages
*/
