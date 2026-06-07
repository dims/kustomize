// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package builtinschema registers the built-in Kubernetes OpenAPI schema with the
// openapi package. Import it (usually for side effects) from any binary or test
// that needs kustomize's built-in-type strategic-merge and namespace-scope
// support:
//
//	import _ "sigs.k8s.io/kustomize/kyaml/openapi/builtinschema"
//
// Importing this package is what links the (large) schema data into a binary.
// Consumers that embed kustomize as a library but do not need the built-in schema
// simply do not import it, and so do not carry the data.
package builtinschema

import (
	"sigs.k8s.io/kustomize/kyaml/openapi"
	"sigs.k8s.io/kustomize/kyaml/openapi/kubernetesapi"
)

func init() {
	openapi.RegisterBuiltinSchema(kubernetesapi.Schema())
}
