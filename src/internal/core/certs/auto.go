package certs

import (
	"crypto/tls"

	"github.com/kodmain/kitsune/src/config"
	"golang.org/x/crypto/acme/autocert"
)

func generateRemoteSignedCert(hosts []string) *tls.Config {
	var certManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Cache:      autocert.DirCache(config.PATH_RUN + "/certs"),
	}

	return certManager.TLSConfig()
}
