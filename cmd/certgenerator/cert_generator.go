// Package cert created for generate a certificate and private key.
package certgenerator

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math"
	"math/big"
	"net"
	"os"
	"time"
)

// GenCert generates a certificate and private key of length lenPrivateKey.
// Return certificate and key files or error.
func GenCert(lenPrivateKey int) (certFile *os.File, privateKeyFile *os.File, err error) {
	keyBytes := make([]byte, 8)
	_, err = rand.Read(keyBytes)
	if err != nil {
		return nil, nil, err
	}
	keyHash := sha1.Sum(keyBytes)
	ski := keyHash[:]

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(math.MaxInt32),
		Subject: pkix.Name{
			Organization: []string{"BGD&Co"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: ski,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, lenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	certB, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certB,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	certFile, err = os.Create("cert.pem")
	if err != nil {
		return nil, nil, err
	}
	defer certFile.Close()

	privateKeyFile, err = os.Create("priv_key.pem")
	if err != nil {
		return nil, nil, err
	}
	defer privateKeyFile.Close()

	_, err = certFile.Write(certPEM.Bytes())
	if err != nil {
		return nil, nil, err
	}

	_, err = privateKeyFile.Write(privateKeyPEM.Bytes())
	if err != nil {
		return nil, nil, err
	}

	return certFile, privateKeyFile, nil
}
