package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

func GenerateCertificate(service, namespace string) ([]byte, []byte, []byte, error) {
	now := time.Now()
	org := []string{
		"postie.chat",
		"h0n9.postie.chat",
		"cloud-secrets.h0n9.postie.chat",
	}
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1994),
		Subject: pkix.Name{
			Organization: org,
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
		return nil, nil, nil, err
	}

	// create CA certificate
	caCertBytes, err := x509.CreateCertificate(
		rand.Reader,
		ca, ca,
		&caPrivKey.PublicKey, caPrivKey,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// encode caCert to PEM format
	caCertPEM := bytes.Buffer{}
	err = pem.Encode(&caCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertBytes,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	cert := &x509.Certificate{
		DNSNames: []string{
			fmt.Sprintf("%s", service),
			fmt.Sprintf("%s.%s", service, namespace),
			fmt.Sprintf("%s.%s.svc", service, namespace),
		},
		SerialNumber: big.NewInt(1994),
		Subject: pkix.Name{
			Organization: org,
			CommonName:   fmt.Sprintf("%s.%s.svc", service, namespace),
		},
		NotBefore:    now,
		NotAfter:     now.AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 9, 9, 4},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage: x509.KeyUsageDigitalSignature,
	}
	serverPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}
	serverPrivKeyBytes := x509.MarshalPKCS1PrivateKey(serverPrivKey)

	// encode serverPrivKey to PEM format
	serverPrivKeyPEM := bytes.Buffer{}
	err = pem.Encode(&serverPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: serverPrivKeyBytes,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// create server certificate
	serverCertBytes, err := x509.CreateCertificate(
		rand.Reader,
		cert, ca,
		&serverPrivKey.PublicKey, caPrivKey,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// encode serverCert to PEM format
	serverCertPEM := bytes.Buffer{}
	err = pem.Encode(&serverCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return caCertPEM.Bytes(), serverCertPEM.Bytes(), serverPrivKeyPEM.Bytes(), nil
}
