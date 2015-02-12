package gpgutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"code.google.com/p/go.crypto/openpgp"
)

type GpgEncrypt struct {
	PublicKeyringPath string
	entityList        []*openpgp.Entity
}

func NewGpgEncrypt(sPublicKeyringPath string) (GpgEncrypt, error) {
	oGpg := GpgEncrypt{}
	oGpg.PublicKeyringPath = sPublicKeyringPath
	err := oGpg.LoadPublicKeyRing()
	return oGpg, err
}

func (g *GpgEncrypt) LoadPublicKeyRing() error {
	keyringFileBuffer, err := os.Open(g.PublicKeyringPath)
	if err != nil {
		return err
	}
	defer keyringFileBuffer.Close()
	entitylist, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err == nil {
		g.entityList = entitylist
	}
	return err
}

func (g *GpgEncrypt) GetKeyByEmail(keyring openpgp.EntityList, email string) *openpgp.Entity {
	for _, entity := range keyring {
		for _, ident := range entity.Identities {
			if ident.UserId.Email == email {
				return entity
			}
		}
	}
	return nil
}

func (g *GpgEncrypt) EncryptStringToFile(plaintext string, sPath string, sEmail string) error {
	rcptPubKey := g.GetKeyByEmail(g.entityList, sEmail)

	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, []*openpgp.Entity{rcptPubKey}, nil, nil, nil)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(plaintext))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(sPath, buf.Bytes(), 0666)
}

func (g *GpgEncrypt) EncryptFile(pathPlain string, pathCrypt string, sEmail string) error {
	bytesPlain, err := ioutil.ReadFile(pathPlain)
	if err != nil {
		return err
	}
	pathCrypt = strings.Trim(pathCrypt, " ")
	if len(pathCrypt) < 1 {
		pathCrypt = pathPlain
	}
	return g.EncryptStringToFile(string(bytesPlain), pathCrypt, sEmail)
}
