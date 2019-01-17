package app

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/gofunct/mamba/runtime"
	"github.com/google/wire"
	"io/ioutil"
)

// applicationSet is the Wire provider set for the Guestbook application that
// does not depend on the underlying platform.
var ApplicationSet = wire.NewSet(
	NewApplication,
	runtime.Set,
)

type Application struct {
	Name    string
	Version string
	Info   string
	*runtime.Service
}

func NewApplication(name string, service *runtime.Service) *Application {
	return &Application{Name: name, Service: service}
}

// LoadECKeyFromFile loads EC key from unencrypted PEM file.
func (a *Application) LoadECKeyFromFile(fileName string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read signing key file: %v", err)
	}

	privateKeyPEM, _ := pem.Decode(privateKeyBytes)
	if privateKeyPEM == nil {
		return nil, fmt.Errorf("failed to decode pem signing key file: %v", err)
	}

	privateKey, err := x509.ParseECPrivateKey(privateKeyPEM.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signing key file: %v", err)
	}

	return privateKey, nil
}