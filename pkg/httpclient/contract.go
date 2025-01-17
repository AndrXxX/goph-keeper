package httpclient

import "crypto/tls"

type tlsConfigProvider interface {
	ForPublicKey() (*tls.Config, error)
}
