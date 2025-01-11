package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type provider struct {
	keyPath string
}

func NewProvider(CryptoKeyPath string) *provider {
	return &provider{keyPath: CryptoKeyPath}
}

func (p *provider) ForPrivateKey() (*tls.Config, error) {
	if p.keyPath == "" {
		return nil, nil
	}
	certPool, err := p.CertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load cert pool: %w", err)
	}
	return &tls.Config{
		ClientCAs: certPool,
	}, nil
}

func (p *provider) ForPublicKey() (*tls.Config, error) {
	if p.keyPath == "" {
		return nil, nil
	}
	certPool, err := p.CertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load cert pool: %w", err)
	}
	return &tls.Config{
		RootCAs: certPool,
	}, nil
}

func (p *provider) CertPool() (*x509.CertPool, error) {
	file, err := os.ReadFile(p.keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read crypto key file: %w", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(file)
	return certPool, nil
}
