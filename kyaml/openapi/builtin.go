// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package openapi

import (
	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// The built-in Kubernetes OpenAPI schema is *registered* by a separate package
// rather than imported here, so that the (large) schema data is linked only into
// binaries that opt in by importing that package. Consumers that embed kustomize
// as a library but do not need the built-in schema therefore do not pull the
// schema data into their builds.
var (
	registeredBuiltinDefs    spec.Definitions
	registeredBuiltinScopes  map[yaml.TypeMeta]bool
	registeredBuiltinVersion string
)

// RegisterBuiltinSchema installs the built-in Kubernetes OpenAPI schema: the
// indexed definitions, the namespace-scope table, and the schema version. It is
// called from an init() in the package that holds the schema data
// (kyaml/openapi/kubernetesapi). If nothing is registered, kustomize runs with no
// built-in schema -- equivalent to SuppressBuiltInSchemaUse.
func RegisterBuiltinSchema(defs spec.Definitions, namespaceScopes map[yaml.TypeMeta]bool, version string) {
	registeredBuiltinDefs = defs
	registeredBuiltinScopes = namespaceScopes
	registeredBuiltinVersion = version
}
