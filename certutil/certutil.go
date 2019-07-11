package certutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
)

// GenerateCSR generates a certificate signing request
func GenerateCSR() (key *rsa.PrivateKey, csrBuf []byte, err error) {
	// Generate
	key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         "",
			Country:            []string{"US"},
			Province:           []string{"CA"},
			Locality:           []string{"San Francisco"},
			Organization:       []string{"symptomatichq"},
			OrganizationalUnit: []string{"."},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return nil, nil, err
	}

	csrBuf = CSRToPem(csr)
	return key, csrBuf, nil
}

// CSRToPem converts a csr into PEM encoded certificate
func CSRToPem(csr []byte) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csr,
		},
	)
}

// CertToPem converts a der encoded certificate (from CreateCertificate) into PEM encoded certificate
func CertToPem(crt []byte) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: crt,
		},
	)
}

// KeyToPem converts an rsa PrivateKey to a PEM encoded private key
func KeyToPem(key *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
}
