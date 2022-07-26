package webhook

import (
	"fmt"
	"strconv"
	"strings"

	csm "github.com/h0n9/cloud-secrets-manager"
)

type Annotations map[string]string

var annotationsAvailable = map[string]bool{
	"provider":  true,
	"secret-id": true,
	"template":  true,
	"output":    true,
	"injected":  true,
}

func ParseAndCheckAnnotations(input Annotations) Annotations {
	output := map[string]string{}
	for key, value := range input {
		subPath := strings.TrimPrefix(key, csm.AnnotationPrefix+"/")
		if subPath == key {
			continue
		}
		if _, exist := annotationsAvailable[subPath]; !exist {
			continue
		}
		output[subPath] = value
	}
	return output
}

func (a Annotations) IsInected() bool {
	value, exist := a["injected"]
	if !exist {
		return false
	}
	injected, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return injected
}

func (a Annotations) getValue(key string) (string, error) {
	value, exist := a[key]
	if !exist {
		return "", fmt.Errorf("failed to read '%s/%s", csm.AnnotationPrefix, key)
	}
	return value, nil
}

func (a Annotations) GetProvider() (string, error) {
	return a.getValue("provider")
}

func (a Annotations) GetSecretID() (string, error) {
	return a.getValue("secret-id")
}

func (a Annotations) GetTemplate() (string, error) {
	return a.getValue("template")
}

func (a Annotations) GetOutput() (string, error) {
	return a.getValue("output")
}
