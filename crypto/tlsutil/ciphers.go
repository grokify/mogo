package tlsutil

import (
	"crypto/tls"

	"github.com/grokify/mogo/type/maputil"
)

func TLS12CiphersStrongMap() map[uint16]string {
	return map[uint16]string{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384: "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		// tls.TLS_DHE_RSA_WITH_AES_256_GCM_SHA384:   "TLS_RSA_WITH_AES_128_CBC_SHA256", // RSA Key Exchange without Forward Secrecy
		// tls.TLS_DHE_RSA_WITH_AES_128_GCM_SHA256:   "TLS_RSA_WITH_RC4_128_SHA",        // RC4-based Ciphers (Insecure due to biases in RC4)
	}
}

// TLS12CiphersWeakMap returns a map of weak TLS 1.2 ciphers.
// Of note, some ciphers are not supported by Go and not included,
// such as `TLS_RSA_WITH_RC4_128_MD5`.
func TLS12CiphersWeakMap() map[uint16]string {
	return map[uint16]string{
		// tls.TLS_RSA_WITH_RC4_128_MD5:                "TLS_RSA_WITH_RC4_128_MD5", // Not available in Go
		// tls.TLS_RSA_WITH_AES_256_CBC_SHA256:         "TLS_RSA_WITH_AES_256_CBC_SHA256", // Not available in Go
		// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384:   "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384", // Not available in Go
		// tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384: "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384", // Not available in Go
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
		tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:        "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",    // RC4-based Ciphers (Insecure due to biases in RC4)
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:     "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA", // 3DES-based Ciphers (Vulnerable to Sweet32 attack)
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:   "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:          "TLS_ECDHE_RSA_WITH_RC4_128_SHA",  // RC4-based Ciphers (Insecure due to biases in RC4)
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:           "TLS_RSA_WITH_3DES_EDE_CBC_SHA",   // 3DES-based Ciphers (Vulnerable to Sweet32 attack)
		tls.TLS_RSA_WITH_AES_128_CBC_SHA:            "TLS_RSA_WITH_AES_128_CBC_SHA",    // RSA Key Exchange without Forward Secrecy
		tls.TLS_RSA_WITH_AES_256_CBC_SHA:            "TLS_RSA_WITH_AES_256_CBC_SHA",    // RSA Key Exchange without Forward Secrecy
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256:         "TLS_RSA_WITH_AES_128_CBC_SHA256", // RSA Key Exchange without Forward Secrecy
		tls.TLS_RSA_WITH_RC4_128_SHA:                "TLS_RSA_WITH_RC4_128_SHA",        // RC4-based Ciphers (Insecure due to biases in RC4)
	}
}

func TLS12CiphersWeak() []string {
	return maputil.ValuesString(TLS12CiphersWeakMap())
}
