package core

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Constants used for the SSL Certificates
const (
	bits          = 2048
	organization  = "Proxify CA"
	country       = "US"
	province      = "CA"
	locality      = "San Francisco"
	streetAddress = "548 Market St"
	postalCode    = "94104"
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

func readCertificationAuthority() (CertificationAuthority, error) {
	caCertLocation := viper.GetString(conf.CaCertLocation)
	caKeyLocation := viper.GetString(conf.CaKeyLocation)
	caKeyPassword := viper.GetString(conf.CaKeyPassword)
	return ReadCertificationAuthority(caCertLocation, caKeyLocation, caKeyPassword)
}

func NewCertSpoofer() (*CertSpoofer, error) {
	ca, err := readCertificationAuthority()
	if err != nil {
		return nil, err
	}
	// Generate one and use it for signing all the generated certs #security :P
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	certSpoofer := CertSpoofer{certAuthority: ca, privateKey: privKey}
	return &certSpoofer, nil
}

type CertSpoofer struct {
	certAuthority CertificationAuthority
	privateKey    *rsa.PrivateKey
}

func (c *CertSpoofer) spoofCertificate(domain string) *x509.Certificate {
	validityYears := viper.GetInt(conf.SpoofedCertValidityYears)

	return &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"You've been pwned buddy!"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Secretville"},
			StreetAddress: []string{"Concealed"},
			PostalCode:    []string{"314159"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(validityYears, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{domain},
	}
}

func (c *CertSpoofer) GenerateSpoofedServerCertificate(domain string) (*tls.Certificate, error) {
	cert := c.spoofCertificate(domain)

	var spoofedCertBytes []byte
	ca := c.certAuthority
	spoofedCertBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &c.privateKey.PublicKey, ca.PrivateKey)
	if err != nil {
		return nil, err
	}

	spoofedCertPem := new(bytes.Buffer)
	err = pem.Encode(spoofedCertPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: spoofedCertBytes,
	})
	if err != nil {
		return nil, err
	}

	spoofedCertPrivKeyPem := new(bytes.Buffer)
	err = pem.Encode(spoofedCertPrivKeyPem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(c.privateKey),
	})
	if err != nil {
		return nil, err
	}

	serverCert, err := tls.X509KeyPair(spoofedCertPem.Bytes(), spoofedCertPrivKeyPem.Bytes())
	return &serverCert, nil

}
