package handler

import (
	"encoding/json"
	"os"
	"text/template"

	"github.com/h0n9/cloud-secrets-manager/provider"
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

func (handler *SecretHandler) Save(secretID, path string) error {
	m, err := handler.Get(secretID)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return handler.template.Execute(file, m)
}
