package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/h0n9/toybox/cloud-secrets-manager/handler"
	"github.com/h0n9/toybox/cloud-secrets-manager/provider"
)

type Mutator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (mutator *Mutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := mutator.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	annotations := ParseAndCheckAnnotations(pod.GetAnnotations())
	if len(annotations) == 0 {
		return admission.Allowed("found no annotations related to cloud-secrets-injector")
	}

	if annotations.IsInected() {
		return admission.Allowed("do nothing as secrets are already injected")
	}

	secretID, err := annotations.GetSecretID()
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	providerStr, err := annotations.GetProvider()
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	tmplStr, err := annotations.GetTemplate()
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	output, err := annotations.GetOutput()
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	tmpl := template.New("cloud-secrets-injector")
	tmpl, err = tmpl.Parse(tmplStr)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	var (
		secretProvider provider.SecretProvider
		secretHandler  *handler.SecretHandler
	)

	switch strings.ToLower(providerStr) {
	case "aws":
		secretProvider, err = provider.NewAWS(ctx)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
	default:
		err = fmt.Errorf("failed to figure out secret provider")
		return admission.Errored(http.StatusBadRequest, err)
	}

	secretHandler, err = handler.NewSecretHandler(secretProvider, tmpl)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	err = secretHandler.Save(secretID, output)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	pod.Annotations[fmt.Sprintf("%s/%s", AnnotationPrefix, "injected")] = "true"
	data, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, data)
}

func (mutator *Mutator) InjectDecoder(decoder *admission.Decoder) error {
	mutator.decoder = decoder
	return nil
}
