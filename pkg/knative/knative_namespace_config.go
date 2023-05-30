/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

        https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package knative

import (
	"context"

	"knative.dev/reconciler-test/pkg/environment"
)

// Deprecated: use environment.WithKnativeNamespace
func WithKnativeNamespace(namespace string) environment.EnvOpts {
	return environment.WithKnativeNamespace(namespace)
}

// Deprecated: use environment.KnativeNamespaceFromContext
func KnativeNamespaceFromContext(ctx context.Context) string {
	return environment.KnativeNamespaceFromContext(ctx)
}
