package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/elazarl/goproxy"
)

// Constants used for the SSL Certificates
const (
	bits          = 2048
	organization  = "Jsubfinder"
	country       = "CA"
	province      = "ON"
	locality      = "Toronto"
	streetAddress = "1 Not Your Business"
	postalCode    = "123456"
)

//https://github.com/Skamaniak/ssl-decryption/blob/cc35498125e66e0ef770a85e4b9b7a5df582de7a/crypto/spoof.go#L65
//https://github.com/Skamaniak/ssl-decryption/blob/cc35498125e66e0ef770a85e4b9b7a5df582de7a/server/web.go#L132
//https://github.com/projectdiscovery/proxify/blob/0fdaa7d0fc4122d1a2e48d054771c744636b5caf/pkg/certs/ca.go

// createCertificateAuthority creates a new certificate authority
func CreateAuthority(certPath, keyPath string) error {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(365*24) * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{organization},
			Country:       []string{country},
			Province:      []string{province},
			Locality:      []string{locality},
			StreetAddress: []string{streetAddress},
			PostalCode:    []string{postalCode},
			CommonName:    organization,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()

	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return err
	}
	return pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
}

func signCertificate(host string) (*tls.Certificate, error) {
	x509ca, err := x509.ParseCertificate(X509pair.Leaf.Raw)
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(365*24) * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Issuer:       x509ca.Subject,
		Subject: pkix.Name{
			Organization:  []string{organization},
			Country:       []string{country},
			Province:      []string{province},
			Locality:      []string{locality},
			StreetAddress: []string{streetAddress},
			PostalCode:    []string{postalCode},
			CommonName:    host,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	certpriv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, x509ca, &certpriv.PublicKey, X509pair.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &tls.Certificate{Certificate: [][]byte{derBytes, X509pair.Leaf.Raw}, PrivateKey: certpriv}, nil
}

func returnCert(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return signCertificate(helloInfo.ServerName)
}

// readCertificateDisk reads a certificate and key file from disk
func ReadCertificateDisk(certFile, keyFile string) error {
	goproxyCa, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0])
	if err != nil {
		return err
	}
	// Check the expiration.
	if time.Now().After(goproxyCa.Leaf.NotAfter) {
		return errors.New("Expired Certificate")
	}

	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}

	return nil
}
