/*
Copyright 2022 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package manifest

import (
	"context"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	pkgsecurity "knative.dev/pkg/test/security"

	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/feature"
)

// PodSecurityCfgFn returns a function for configuring security context for Pod, depending
// on security settings of the enclosing namespace.
func PodSecurityCfgFn(ctx context.Context, t feature.T) CfgFn {
	namespace := environment.FromContext(ctx).Namespace()
	restrictedMode, err := pkgsecurity.IsRestrictedPodSecurityEnforced(ctx, kubeclient.Get(ctx), namespace)
	if err != nil {
		t.Fatalf("Error while checking restricted pod security mode for namespace %s", namespace)
	}
	if restrictedMode {
		return WithDefaultPodSecurityContext
	}
	return func(map[string]interface{}) {}
}

// WithAnnotations returns a function for configuring annototations of the resource
func WithAnnotations(annotations map[string]interface{}) CfgFn {
	return func(cfg map[string]interface{}) {
		if original, ok := cfg["annotations"]; ok {
			appendToOriginal(original, annotations)
			return
		}
		cfg["annotations"] = annotations
	}
}

// WithPodAnnotations appends pod annotations (usually used by types where pod template is embedded)
func WithPodAnnotations(additional map[string]interface{}) CfgFn {
	return func(cfg map[string]interface{}) {
		if ann, ok := cfg["podannotations"]; ok {
			appendToOriginal(ann, additional)
			return
		}
		cfg["podannotations"] = additional
	}
}

func appendToOriginal(original interface{}, additional map[string]interface{}) {
	annotations := original.(map[string]interface{})
	for k, v := range additional {
		// Only add the unspecified ones
		if _, ok := annotations[k]; !ok {
			annotations[k] = v
		}
	}
}

// WithLabels returns a function for configuring labels of the resource
func WithLabels(labels map[string]string) CfgFn {
	return func(cfg map[string]interface{}) {
		if labels != nil {
			cfg["labels"] = labels
		}
	}
}

func WithIstioPodAnnotations(cfg map[string]interface{}) {
	podAnnotations := map[string]interface{}{
		"sidecar.istio.io/inject":                "true",
		"sidecar.istio.io/rewriteAppHTTPProbers": "true",
	}

	WithAnnotations(podAnnotations)(cfg)
	WithPodAnnotations(podAnnotations)(cfg)
}

func WithDefaultPodSecurityContext(cfg map[string]interface{}) {
	if _, set := cfg["podSecurityContext"]; !set {
		cfg["podSecurityContext"] = map[string]interface{}{}
	}
	podSecurityContext := cfg["podSecurityContext"].(map[string]interface{})
	podSecurityContext["runAsNonRoot"] = pkgsecurity.DefaultPodSecurityContext.RunAsNonRoot
	podSecurityContext["seccompProfile"] = map[string]interface{}{}
	seccompProfile := podSecurityContext["seccompProfile"].(map[string]interface{})
	seccompProfile["type"] = pkgsecurity.DefaultPodSecurityContext.SeccompProfile.Type

	if _, set := cfg["containerSecurityContext"]; !set {
		cfg["containerSecurityContext"] = map[string]interface{}{}
	}
	containerSecurityContext := cfg["containerSecurityContext"].(map[string]interface{})
	containerSecurityContext["allowPrivilegeEscalation"] =
		pkgsecurity.DefaultContainerSecurityContext.AllowPrivilegeEscalation
	containerSecurityContext["capabilities"] = map[string]interface{}{}
	capabilities := containerSecurityContext["capabilities"].(map[string]interface{})
	if len(pkgsecurity.DefaultContainerSecurityContext.Capabilities.Drop) != 0 {
		capabilities["drop"] = []string{}
		for _, drop := range pkgsecurity.DefaultContainerSecurityContext.Capabilities.Drop {
			capabilities["drop"] = append(capabilities["drop"].([]string), string(drop))
		}
	}
	if len(pkgsecurity.DefaultContainerSecurityContext.Capabilities.Add) != 0 {
		capabilities["add"] = []string{}
		for _, drop := range pkgsecurity.DefaultContainerSecurityContext.Capabilities.Drop {
			capabilities["add"] = append(capabilities["add"].([]string), string(drop))
		}
	}
}
