package webhook

import (
	"fmt"
	"strconv"
	"strings"

	csm "github.com/h0n9/cloud-secrets-manager"
)

type AnnotationSet map[string]Annotations

var AnnotationMap = map[string]string{
	"cloud-secrets-manager.h0n9.postie.chat/provider":  "provider",
	"cloud-secrets-manager.h0n9.postie.chat/secret-id": "secret-id",
	"cloud-secrets-manager.h0n9.postie.chat/template":  "template",
	"cloud-secrets-manager.h0n9.postie.chat/output":    "output",
	"cloud-secrets-manager.h0n9.postie.chat/injected":  "injected",
}

func ParseAnnotationSet(input map[string]string) AnnotationSet {
	var (
		output               = AnnotationSet{}
		exist                bool
		subPath, full, short string
	)
	// TODO: enhance O(5N)
	for key, value := range input {
		for full, short = range AnnotationMap {
			subPath = strings.TrimPrefix(key, full)
			if subPath == key {
				continue
			}
			subPath = strings.TrimPrefix(subPath, "-")
			if _, exist = output[subPath]; !exist {
				output[subPath] = Annotations{}
			}
			output[subPath][short] = value
			break
		}
	}
	return output
}

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
