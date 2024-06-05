package handler

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"

	"github.com/h0n9/cloud-secrets-manager/provider"
	"github.com/h0n9/cloud-secrets-manager/util"
)

type SecretHandlerFunc func(string) (string, error)

type SecretHandler struct {
	provider provider.SecretProvider
	template *template.Template
}

func NewSecretHandler(provider provider.SecretProvider, tmpl *template.Template) (*SecretHandler, error) {
	return &SecretHandler{provider: provider, template: tmpl}, nil
}

func (handler *SecretHandler) Get(secretID string) (map[string]interface{}, error) {
	secretValue, err := handler.provider.GetSecretValue(secretID)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(secretValue), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (handler *SecretHandler) Save(secretID, path string, decodeBase64EncodedSecret bool) error {
	// get secret
	m, err := handler.Get(secretID)
	if err != nil {
		return err
	}

	// create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// if secret is not base64 encoded, write it to file and return
	if !decodeBase64EncodedSecret {
		return handler.template.Execute(file, m)
	}

	var (
		buff          bytes.Buffer
		decodedSecret []byte
	)

	// execute template
	err = handler.template.Execute(&buff, m)
	if err != nil {
		return err
	}

	// decode base64 encoded secret
	decodedSecret, err = util.DecodeBase64StrToBytes(buff.String())
	if err != nil {
		return err
	}

	// write decoded secret to file
	_, err = file.Write(decodedSecret)
	return err
}
