package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	l "github.com/hiddengearz/jsubfinder/core/logger"
	"github.com/mitchellh/go-homedir"
)

func GenerateCert() (err error) {

	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	sshFolder := home + "/.ssh/"
	if !folderExists(sshFolder) {
		l.Log.Debug("Folder " + sshFolder + " doesnt exist. Please create it")
		return errors.New("Folder " + sshFolder + " doesnt exist. Please create it")
	}

	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		l.Log.Debug("Cannot generate RSA key\n")
		return
	}
	publickey := &privatekey.PublicKey

	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create(sshFolder + "jsubfinder")
	if err != nil {
		l.Log.Debug("error when create private.pem: %s \n", err)
		return
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		l.Log.Debug("error when encode private pem: %s \n", err)
		return
	}

	// dump public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		l.Log.Debug("error when dumping publickey: %s \n", err)
		return
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create(sshFolder + "jsubfinder.pub")
	if err != nil {
		l.Log.Debug("error when create public.pem: %s \n", err)
		return
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		l.Log.Debug("error when encode public pem: %s \n", err)
		return
	}

	return nil

}

func ParseCertAndKey(certPEM, keyPEM []byte) (*x509.Certificate, *rsa.PrivateKey, error) {
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, nil, err
	}

	key, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("private key with unexpected type %T", key)
	}
	return cert, key, nil
}
