package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	csm "github.com/h0n9/cloud-secrets-manager"
	"github.com/h0n9/cloud-secrets-manager/util"
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

	annotationSet := ParseAnnotationSet(pod.GetAnnotations())
	if len(annotationSet) == 0 {
		return admission.Allowed("found no annotations related to cloud-secrets-injector")
	}

	for secretName, annotations := range annotationSet {
		if annotations.IsInected() {
			// return admission.Allowed("do nothing as secrets are already injected")
			continue
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

		// prepare injector name for general use
		injectorName := "cloud-secrets-injector"
		if secretName != "" {
			injectorName = injectorName + "-" + secretName
		}

		// append volume to pod volumes
		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name:         injectorName,
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		})

		// inject sidecar
		pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
			Name:  injectorName,
			Image: mutator.InjectorImage,
			Args: []string{
				"injector",
				"run",
				fmt.Sprintf("--provider=%s", providerStr),
				fmt.Sprintf("--secret-id=%s", secretID),
				fmt.Sprintf("--template=%s", util.EncodeBase64(tmplStr)),
				fmt.Sprintf("--output=%s", output),
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      injectorName,
					MountPath: filepath.Dir(output),
				},
			},
		})

		// mount volume to every containers
		for i := range pod.Spec.Containers {
			pod.Spec.Containers[i].VolumeMounts = append(pod.Spec.Containers[i].VolumeMounts, corev1.VolumeMount{
				Name:      injectorName,
				MountPath: filepath.Dir(output),
				SubPath:   filepath.Base(output),
			})
		}

		// set annotation for injection to true
		injected := "injected"
		if secretName != "" {
			injected = injected + "-" + secretName
		}
		pod.Annotations[fmt.Sprintf("%s/%s", csm.AnnotationPrefix, injected)] = "true"
	}

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
