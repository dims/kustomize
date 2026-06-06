// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package kubernetesapi

// BuiltinSchema is the compact, purpose-built form of the built-in Kubernetes
// OpenAPI schema. It holds only what kustomize reads: the structural type graph
// (enough to resolve field and array-element types), the strategic-merge and
// group-version-kind extensions, and a precomputed namespace-scope table.
//
// It is produced by kubernetesapi/scripts/extractschema from a full swagger.pb,
// embedded per version, and rebuilt into kube-openapi spec.Schema objects at
// load time by the openapi package.
type BuiltinSchema struct {
	// Definitions maps an OpenAPI definition name to its trimmed schema.
	Definitions map[string]*SchemaNode `json:"definitions"`
	// Scopes records namespace- vs cluster-scope per resource type. Its set of
	// types is derived from API paths and differs from the Definitions' GVKs.
	Scopes []ScopeEntry `json:"scopes"`
}

// SchemaNode is the subset of an OpenAPI schema that kustomize navigates and
// merges against. Documentation, validation constraints, defaults, examples and
// everything else kustomize never reads are omitted.
type SchemaNode struct {
	Type                 []string               `json:"type,omitempty"`
	Ref                  string                 `json:"ref,omitempty"` // definition name, without the "#/definitions/" prefix
	Items                *SchemaNode            `json:"items,omitempty"`
	Properties           map[string]*SchemaNode `json:"properties,omitempty"`
	AdditionalProperties *SchemaNode            `json:"additionalProperties,omitempty"` // value schema of a map field
	PatchStrategy        string                 `json:"patchStrategy,omitempty"`        // x-kubernetes-patch-strategy
	PatchMergeKey        string                 `json:"patchMergeKey,omitempty"`        // x-kubernetes-patch-merge-key
	ListMapKeys          []string               `json:"listMapKeys,omitempty"`          // x-kubernetes-list-map-keys
	GroupVersionKinds    []GVK                  `json:"groupVersionKinds,omitempty"`    // x-kubernetes-group-version-kind
}

// GVK is a group/version/kind tuple.
type GVK struct {
	Group   string `json:"group,omitempty"`
	Version string `json:"version,omitempty"`
	Kind    string `json:"kind,omitempty"`
}

// ScopeEntry records whether a resource type is namespace-scoped.
type ScopeEntry struct {
	Group      string `json:"group,omitempty"`
	Version    string `json:"version,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Namespaced bool   `json:"namespaced,omitempty"`
}
