package webhook

import (
	"fmt"
	"strconv"
	"strings"

	csm "github.com/h0n9/cloud-secrets-manager"
)

const (
	AnnotationProvider = "provider"
	AnnotationSecretID = "secret-id"
	AnnotationTemplate = "template"
	AnnotationOutput   = "output"
	AnnotationInjected = "injected"
)

type AnnotationSet map[string]Annotations

var AnnotationMap = map[string]string{
	fmt.Sprintf("%s/%s", csm.AnnotationPrefix, AnnotationProvider): AnnotationProvider,
	fmt.Sprintf("%s/%s", csm.AnnotationPrefix, AnnotationSecretID): AnnotationSecretID,
	fmt.Sprintf("%s/%s", csm.AnnotationPrefix, AnnotationTemplate): AnnotationTemplate,
	fmt.Sprintf("%s/%s", csm.AnnotationPrefix, AnnotationOutput):   AnnotationOutput,
	fmt.Sprintf("%s/%s", csm.AnnotationPrefix, AnnotationInjected): AnnotationInjected,
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
	AnnotationProvider: true,
	AnnotationSecretID: true,
	AnnotationTemplate: true,
	AnnotationOutput:   true,
	AnnotationInjected: true,
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
	value, exist := a[AnnotationInjected]
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
	return a.getValue(AnnotationProvider)
}

func (a Annotations) GetSecretID() (string, error) {
	return a.getValue(AnnotationSecretID)
}

func (a Annotations) GetTemplate() (string, error) {
	return a.getValue(AnnotationTemplate)
}

func (a Annotations) GetOutput() (string, error) {
	return a.getValue(AnnotationOutput)
}
