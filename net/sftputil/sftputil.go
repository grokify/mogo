package sftputil

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	DefaultPort = int64(22)
)

type SftpClient struct {
	Client     *sftp.Client
	hostname   string
	port       int64
	username   string
	keyPath    string
	password   string
	authMethod string
}

func NewSftpClientWithPassword(hostname string, port int64, username string, password string) (SftpClient, error) {
	oSftp := SftpClient{}
	oSftp.authMethod = "password"
	oSftp.SetCredentials(hostname, port, username, "", password)
	err := oSftp.Connect()
	return oSftp, err
}

func NewSftpClientWithKeyPath(hostname string, port int64, username string, keyPath string) (SftpClient, error) {
	oSftp := SftpClient{}
	oSftp.authMethod = "key"
	oSftp.SetCredentials(hostname, port, username, keyPath, "")
	err := oSftp.Connect()
	return oSftp, err
}

func (c *SftpClient) SetCredentials(hostname string, port int64, username string, keyPath string, password string) {
	c.hostname = hostname
	c.port = port
	c.username = username
	c.keyPath = keyPath
	c.password = password
}

func (c *SftpClient) GetKey(sKeyPath string) (key ssh.Signer, err error) {
	buf, err := ioutil.ReadFile(sKeyPath)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(buf)
}

func (c *SftpClient) Connect() error {
	auth := []ssh.AuthMethod{}
	if c.authMethod == "key" {
		key, _ := c.GetKey(c.keyPath)

		auth = []ssh.AuthMethod{
			ssh.PublicKeys(key),
		}
	} else if c.authMethod == "password" {
		auth = []ssh.AuthMethod{
			ssh.Password(c.password),
		}
	}
	config := &ssh.ClientConfig{
		User: c.username,
		Auth: auth,
	}
	sHost := strings.Join([]string{c.hostname, strconv.FormatInt(c.port, 10)}, ":")

	sshClient, err := ssh.Dial("tcp", sHost, config)
	if err != nil {
		return err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err == nil {
		c.Client = sftpClient
	}
	return err
}

func (c *SftpClient) Get(sPathRem string, sPathLoc string) error {
	fi, err := c.Client.Open(sPathRem)
	if err != nil {
		rex := regexp.MustCompile(`\(SSH_FX_FAILURE\)`) // Hack, SFTP failure at 100
		res := rex.FindStringSubmatch(err.Error())
		if len(res) > 0 {
			c.Connect()
			return c.Get(sPathRem, sPathLoc)
		}
		log.Fatal(err)
	}
	fo, err := os.Create(sPathLoc)
	if err != nil {
		return err
	}
	defer fo.Close()
	_, err = io.Copy(fo, fi)
	return err
}

func (c *SftpClient) Put(sPathLoc string, sPathRem string) error {
	ab, err := ioutil.ReadFile(sPathLoc)
	if err != nil {
		return err
	}
	fRem, err := c.Client.Create(sPathRem)
	if err == nil {
		fRem.Write(ab)
	}
	return err
}
