package certs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
)

// generateSelfSignedCert génère un certificat auto-signé et renvoie un tls.Config.
func generateSelfSignedCert(hosts []string) *tls.Config {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if logger.Fatal(err) {
		return nil
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if logger.Fatal(err) {
		return nil
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Kitsune"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              hosts,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if logger.Fatal(err) {
		return nil
	}

	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	return tlsConfig
}
