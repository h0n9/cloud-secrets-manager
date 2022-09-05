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
	// error cases
	parsed, err := ParseAndCheckAnnotations(Annotations{
		"cloud-secrets-manager.h0n9.postie.chat/volume-path": "/envs", // ❌: unsupported
	})
	assert.Error(t, err)
	assert.EqualValues(t, Annotations{}, parsed)
	parsed, err = ParseAndCheckAnnotations(Annotations{
		"cloud-secrets-manager.h0n9.postie.chat/secret": "my-precious-secret", // ❌ secret 💔 output
		"cloud-secrets-manager.h0n9.postie.chat/output": "envs",               // ❌ secret 💔 output
	})
	assert.Error(t, err)
	assert.EqualValues(t, Annotations{}, parsed)
	parsed, err = ParseAndCheckAnnotations(Annotations{
		"cloud-secrets-manager.h0n9.postie.chat/secret":   "my-precious-secret", // ❌ secret 💔 template
		"cloud-secrets-manager.h0n9.postie.chat/template": SampleTemplate,       // ❌ secret 💔 template
	})
	assert.Error(t, err)
	assert.EqualValues(t, Annotations{}, parsed)

	// ignore cases
	parsed, err = ParseAndCheckAnnotations(Annotations{
		"vault.hashicorp.com/secret-volume-path-SECRET-NAME-foobar": "/envs", // ❌: non related annotation
	})
	assert.NoError(t, err)
	assert.EqualValues(t, Annotations{}, parsed)
	parsed, err = ParseAndCheckAnnotations(Annotations{
		"cloud-secrets-manager.h0n9.postie.chat": "h0n9", // ❌: non subpath
	})
	assert.NoError(t, err)
	assert.EqualValues(t, Annotations{}, parsed)
	parsed, err = ParseAndCheckAnnotations(Annotations{
		"cloud-secrets-manager.h0n9.posite.chat/template": SampleTemplate, // ❌: typo
	})
	assert.NoError(t, err)
	assert.EqualValues(t, Annotations{}, parsed)

	// good cases
	parsed, err = ParseAndCheckAnnotations(Annotations{
		"cloud-secrets-manager.h0n9.postie.chat/provider":  "aws",               // ✅
		"cloud-secrets-manager.h0n9.postie.chat/secret-id": "life-is-beautiful", // ✅
		"cloud-secrets-manager.h0n9.postie.chat/output":    "/envs",             // ✅
		"cloud-secrets-manager.h0n9.postie.chat/template":  SampleTemplate,      // ✅
		"cloud-secrets-manager.h0n9.postie.chat/injected":  "true",              // ✅
	})
	assert.NoError(t, err)
	assert.EqualValues(t, Annotations{
		"provider":  "aws",
		"secret-id": "life-is-beautiful",
		"template":  SampleTemplate,
		"output":    "/envs",
		"injected":  "true",
	}, parsed)
	parsed, err = ParseAndCheckAnnotations(Annotations{
		"cloud-secrets-manager.h0n9.postie.chat/provider":  "aws",                // ✅
		"cloud-secrets-manager.h0n9.postie.chat/secret-id": "life-is-beautiful",  // ✅
		"cloud-secrets-manager.h0n9.postie.chat/secret":    "my-precious-secret", // ✅
		"cloud-secrets-manager.h0n9.postie.chat/injected":  "true",               // ✅
	})
	assert.NoError(t, err)
	assert.EqualValues(t, Annotations{
		"provider":  "aws",
		"secret-id": "life-is-beautiful",
		"secret":    "my-precious-secret",
		"injected":  "true",
	}, parsed)
	parsed, err = ParseAndCheckAnnotations(Annotations{}) // ✅ empty is good
	assert.NoError(t, err)
	assert.EqualValues(t, Annotations{}, parsed)
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
