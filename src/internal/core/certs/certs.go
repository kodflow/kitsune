package certs

import (
	"crypto/tls"
	"fmt"
)

func TLSConfigFor(domain string, subs ...string) *tls.Config {
	if domain == "" {
		domain = "localhost"
	}

	hosts := generateHosts(domain, subs...)

	if domain == "localhost" {
		return generateSelfSignedCert(hosts)
	}

	return generateRemoteSignedCert(hosts)
}

func generateHosts(domain string, subs ...string) []string {
	var combinations []string
	combinations = append(combinations, domain)
	for _, sub := range subs {
		combinations = append(combinations, fmt.Sprintf("%s.%s", sub, domain))
	}

	return combinations
}
