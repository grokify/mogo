package gpgutil

import (
	"bytes"
	"os"
	"strings"

	// "golang.org/x/crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp"
)

type GPGEncrypt struct {
	PublicKeyringPath string
	entityList        []*openpgp.Entity
}

func NewGPGEncrypt(sPublicKeyringPath string) (GPGEncrypt, error) {
	oGpg := GPGEncrypt{}
	oGpg.PublicKeyringPath = sPublicKeyringPath
	err := oGpg.LoadPublicKeyRing()
	return oGpg, err
}

func (g *GPGEncrypt) LoadPublicKeyRing() error {
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

func (g *GPGEncrypt) GetKeyByEmail(keyring openpgp.EntityList, email string) *openpgp.Entity {
	for _, entity := range keyring {
		for _, ident := range entity.Identities {
			if ident.UserId.Email == email {
				return entity
			}
		}
	}
	return nil
}

func (g *GPGEncrypt) EncryptStringToFile(plaintext string, sPath string, sEmail string) error {
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
	return os.WriteFile(sPath, buf.Bytes(), 0600)
}

func (g *GPGEncrypt) EncryptFile(pathPlain string, pathCrypt string, sEmail string) error {
	bytesPlain, err := os.ReadFile(pathPlain)
	if err != nil {
		return err
	}
	pathCrypt = strings.Trim(pathCrypt, " ")
	if len(pathCrypt) < 1 {
		pathCrypt = pathPlain
	}
	return g.EncryptStringToFile(string(bytesPlain), pathCrypt, sEmail)
}
