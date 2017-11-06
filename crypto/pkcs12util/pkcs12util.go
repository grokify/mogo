package pkcs12util

import (
	"fmt"
	"strings"
)

const (
	// > openssl pkcs12 -export -in cert.pem -inkey "privateKey.pem" -certfile cert.pem -out myProject_keyAndCertBundle.p12
	OpenSSLCreateFormat = `openssl pkcs12 -export -in %v -inkey %v -certfile %v -out %v`
)

type Options struct {
	// use certificate.crt as the certificate the private key will be combined with.
	In string
	// use the private key file privateKey.key as the private key to combine with the certificate.
	InKey string
	// This is optional, this is if you have any additional certificates you would like to include in the PFX file.
	CertFile string
	// export and save the PFX file as certificate.pfx
	Out string
}

func (opts *Options) TrimSpace() {
	opts.In = strings.TrimSpace(opts.In)
	opts.InKey = strings.TrimSpace(opts.InKey)
	opts.CertFile = strings.TrimSpace(opts.CertFile)
	opts.Out = strings.TrimSpace(opts.Out)
}

// CreateP12File creates a PKCS 12/PFX file
// https://www.ssl.com/how-to/create-a-pfx-p12-certificate-file-using-openssl/
func (opts *Options) CreateCommand() string {
	parts := []string{"openssl pkcs12 -export"}
	return opts.Stringify(parts)
}

func (opts *Options) InfoCommand() string {
	parts := []string{"openssl pkcs12 -info"}
	return opts.Stringify(parts)
}

func (opts *Options) Stringify(parts []string) string {
	opts.TrimSpace()
	if len(opts.Out) > 0 {
		parts = append(parts, fmt.Sprintf("-out %v", opts.Out))
	}
	if len(opts.InKey) > 0 {
		parts = append(parts, fmt.Sprintf("-inkey %v", opts.InKey))
	}
	if len(opts.In) > 0 {
		parts = append(parts, fmt.Sprintf("-in %v", opts.In))
	}
	if len(opts.CertFile) > 0 {
		parts = append(parts, fmt.Sprintf("-in %v", opts.CertFile))
	}
	return strings.Join(parts, " ")
}

//openssl pkcs12 -info -in

/*
// CreateP12File creates a PKCS 12/PFX file
// https://www.ssl.com/how-to/create-a-pfx-p12-certificate-file-using-openssl/
func CreateP12File(opts CreateOptions) error {
	_ = opts.CreateCommand()
	return nil
}
*/
