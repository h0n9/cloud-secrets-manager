package webhook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SampleTemplate = `
{{ range $k, $v := . }}export {{ $k }}={{ $v }}
{{ end }}
	`
)

func TestParseAndCheckAnnotations(t *testing.T) {
	input := map[string]string{
		"cloud-secrets-manager.h0n9.postie.chat/provider":           "aws",               // ✅
		"cloud-secrets-manager.h0n9.postie.chat/secret-id":          "life-is-beautiful", // ✅
		"cloud-secrets-manager.h0n9.postie.chat/output":             "/envs",             // ✅
		"cloud-secrets-manager.h0n9.postie.chat/template":           SampleTemplate,      // ✅
		"cloud-secrets-manager.h0n9.postie.chat/injected":           "true",              // ✅
		"cloud-secrets-manager.h0n9.postie.chat/decode-base64":      "true",              // ✅
		"cloud-secrets-manager.h0n9.posite.chat/template":           SampleTemplate,      // ❌: typo
		"cloud-secrets-manager.h0n9.postie.chat/volume-path":        "/envs",             // ❌: unsupported
		"cloud-secrets-manager.h0n9.postie.chat":                    "h0n9",              // ❌: non subpath
		"vault.hashicorp.com/secret-volume-path-SECRET-NAME-foobar": "/envs",             // ❌: non related annotation
	}
	output := ParseAndCheckAnnotations(input)
	expectedOutput := Annotations{
		"provider":      "aws",
		"secret-id":     "life-is-beautiful",
		"template":      SampleTemplate,
		"output":        "/envs",
		"decode-base64": "true",
		"injected":      "true",
	}
	assert.EqualValues(t, expectedOutput, output)
}

func TestParseAnnotationSet(t *testing.T) {
	input := map[string]string{
		"cloud-secrets-manager.h0n9.postie.chat/provider":           "aws",               // ✅
		"cloud-secrets-manager.h0n9.postie.chat/secret-id":          "life-is-beautiful", // ✅
		"cloud-secrets-manager.h0n9.postie.chat/output":             "/envs",             // ✅
		"cloud-secrets-manager.h0n9.postie.chat/template":           SampleTemplate,      // ✅
		"cloud-secrets-manager.h0n9.postie.chat/injected":           "true",              // ✅
		"cloud-secrets-manager.h0n9.postie.chat/decode-base64":      "true",              // ✅
		"cloud-secrets-manager.h0n9.posite.chat/template":           SampleTemplate,      // ❌: typo
		"cloud-secrets-manager.h0n9.postie.chat/volume-path":        "/envs",             // ❌: unsupported
		"cloud-secrets-manager.h0n9.postie.chat":                    "h0n9",              // ❌: non subpath
		"vault.hashicorp.com/secret-volume-path-SECRET-NAME-foobar": "/envs",             // ❌: non related annotation
		"cloud-secrets-manager.h0n9.postie.chat/provider-test":      "gcp",               // ✅
		"cloud-secrets-manager.h0n9.postie.chat/secret-id-test":     "life-is-wonderful", // ✅
		"cloud-secrets-manager.h0n9.postie.chat/output-test":        "/secrets/envs",     // ✅
		"cloud-secrets-manager.h0n9.postie.chat/template-test":      SampleTemplate,      // ✅
		"cloud-secrets-manager.h0n9.postie.chat/injected-test":      "false",             // ✅
		"cloud-secrets-manager.h0n9.postie.chat/decode-base64-test": "false",             // ✅
	}
	output := ParseAnnotationSet(input)
	expectedOutput := AnnotationSet{
		"": Annotations{
			"provider":      "aws",
			"secret-id":     "life-is-beautiful",
			"template":      SampleTemplate,
			"output":        "/envs",
			"decode-base64": "true",
			"injected":      "true",
		},
		"test": Annotations{
			"provider":      "gcp",
			"secret-id":     "life-is-wonderful",
			"template":      SampleTemplate,
			"output":        "/secrets/envs",
			"decode-base64": "false",
			"injected":      "false",
		},
	}
	assert.EqualValues(t, expectedOutput, output)
}

func TestAnnotationsIsInjected(t *testing.T) {
	annotations := Annotations{}
	assert.False(t, annotations.IsInected())
	annotations = Annotations{"injected": "false"}
	assert.False(t, annotations.IsInected())
	annotations = Annotations{"injected": "x"}
	assert.False(t, annotations.IsInected())
	annotations = Annotations{"injected": "ture"}
	assert.False(t, annotations.IsInected())
	annotations = Annotations{"injected": "t"}
	assert.True(t, annotations.IsInected())
	annotations = Annotations{"injected": "true"}
	assert.True(t, annotations.IsInected())
}
