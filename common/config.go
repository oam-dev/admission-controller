package common

import (
	"crypto/tls"
	"flag"

	"k8s.io/klog"
)

// Config contains the server (the webhook) cert and key.
type Config struct {
	CertFile string
	KeyFile  string
}

func (c *Config) AddFlags() {
	flag.StringVar(&c.CertFile, "tls-cert-file", "/etc/certs/cert.pem", ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	flag.StringVar(&c.KeyFile, "tls-private-key-file", "/etc/certs/key.pem", ""+
		"File containing the default x509 private key matching --tls-cert-file.")
}

func ConfigTLS(config Config) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		klog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls with apiserver cert
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}
