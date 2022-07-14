package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/h0n9/toybox/cloud-secrets-manager/util"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type Mutator struct {
	Client        client.Client
	InjectorImage string
	decoder       *admission.Decoder
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

	// append 'cloud-secrets-injector' volume to pod volumes
	volumeName := "cloud-secrets-injector"
	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name:         volumeName,
		VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
	})

	// inject sidecar
	pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
		Name:  "cloud-secrets-injector",
		Image: mutator.InjectorImage,
		Env: []corev1.EnvVar{
			{Name: "PROVIDER_NAME", Value: providerStr},
			{Name: "SECRET_ID", Value: secretID},
			{Name: "TEMPLATE_BASE64", Value: util.EncodeBase64(tmplStr)},
			{Name: "OUTPUT_FILE", Value: output},
		},
		Args: []string{"injector", "run"},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      volumeName,
				MountPath: filepath.Dir(output),
			},
		},
	})

	// mount volume to every containers
	for _, container := range pod.Spec.Containers {
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      volumeName,
			MountPath: filepath.Dir(output),
		})
	}

	// set annotation for injection to true
	pod.Annotations[fmt.Sprintf("%s/%s", AnnotationPrefix, "injected")] = "true"

	// marshal pod struct into bytes slice
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
