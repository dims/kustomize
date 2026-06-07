// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package kubernetesapi

import (
	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Schema rebuilds the compact built-in schema into the form the openapi package
// consumes: indexed kube-openapi definitions, a namespace-scope table, and the
// schema version.
//
// This package intentionally does NOT import the openapi package or register
// anything on import: registration is done by the separate builtinschema package,
// so the (large) schema data is linked only into binaries that opt in by
// importing that package. Keeping this package free of an openapi import also
// lets openapi's own tests obtain the schema without an import cycle.
func Schema() (defs spec.Definitions, namespaceScopes map[yaml.TypeMeta]bool, version string) {
	defs = spec.Definitions{}
	for name, node := range Builtin.Definitions {
		defs[name] = rebuildSchema(node)
	}
	namespaceScopes = make(map[yaml.TypeMeta]bool, len(Builtin.Scopes))
	for _, s := range Builtin.Scopes {
		apiVersion := s.Version
		if s.Group != "" {
			apiVersion = s.Group + "/" + s.Version
		}
		namespaceScopes[yaml.TypeMeta{APIVersion: apiVersion, Kind: s.Kind}] = s.Namespaced
	}
	return defs, namespaceScopes, DefaultOpenAPI
}

// rebuildSchema reconstructs a kube-openapi spec.Schema from a compact node,
// reproducing the fields and extension value types kustomize reads: the
// x-kubernetes-* extensions are rebuilt with the dynamic types the openapi
// ResourceSchema accessors assert on (string, []interface{}, and []interface{}
// of map[string]interface{}).
func rebuildSchema(n *SchemaNode) spec.Schema {
	var s spec.Schema
	if n == nil {
		return s
	}
	if len(n.Type) > 0 {
		s.Type = spec.StringOrArray(n.Type)
	}
	if n.Ref != "" {
		s.Ref = spec.MustCreateRef("#/definitions/" + n.Ref)
	}
	if n.Items != nil {
		item := rebuildSchema(n.Items)
		s.Items = &spec.SchemaOrArray{Schema: &item}
	}
	if len(n.Properties) > 0 {
		s.Properties = make(map[string]spec.Schema, len(n.Properties))
		for k, v := range n.Properties {
			s.Properties[k] = rebuildSchema(v)
		}
	}
	if n.AdditionalProperties != nil {
		ap := rebuildSchema(n.AdditionalProperties)
		s.AdditionalProperties = &spec.SchemaOrBool{Allows: true, Schema: &ap}
	}

	ext := spec.Extensions{}
	if n.PatchStrategy != "" {
		ext["x-kubernetes-patch-strategy"] = n.PatchStrategy
	}
	if n.PatchMergeKey != "" {
		ext["x-kubernetes-patch-merge-key"] = n.PatchMergeKey
	}
	if len(n.ListMapKeys) > 0 {
		keys := make([]interface{}, len(n.ListMapKeys))
		for i, k := range n.ListMapKeys {
			keys[i] = k
		}
		ext["x-kubernetes-list-map-keys"] = keys
	}
	if len(n.GroupVersionKinds) > 0 {
		gvks := make([]interface{}, len(n.GroupVersionKinds))
		for i, g := range n.GroupVersionKinds {
			gvks[i] = map[string]interface{}{"group": g.Group, "version": g.Version, "kind": g.Kind}
		}
		ext["x-kubernetes-group-version-kind"] = gvks
	}
	if len(ext) > 0 {
		s.Extensions = ext
	}
	return s
}
