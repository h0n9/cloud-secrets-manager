package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func GenerateCACertificate() ([]byte, error) {
	now := time.Now()
	ca := x509.Certificate{
		SerialNumber: big.NewInt(1994),
		Subject: pkix.Name{
			Organization: []string{
				"postie.chat",
				"h0n9.postie.chat",
				"cloud-secrets.h0n9.postie.chat",
			},
		},
		NotBefore: now,
		NotAfter:  now.AddDate(10, 0, 0),
		IsCA:      true,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	caBytes, err := x509.CreateCertificate(
		rand.Reader,
		&ca, &ca,
		caPrivKey.Public(), caPrivKey,
	)
	if err != nil {
		return nil, err
	}
	return caBytes, nil
}

func GenerateCertificate() ([]byte, error) {
	return nil, nil
}
