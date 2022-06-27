package webhook

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type Validator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (validator *Validator) Handle(ctx context.Context, req admission.Request) admission.Response {
	return admission.Response{}
}

func (validator *Validator) InjectDecoder(decoder *admission.Decoder) error {
	validator.decoder = decoder
	return nil
}
