// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package openapi

import (
	"encoding/json"

	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/openapi/kubernetesapi"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// parseCompactBuiltinSchema rebuilds the global schema from the compact,
// purpose-built built-in schema produced by kubernetesapi/scripts/extractschema.
// It reconstructs the kube-openapi spec.Schema objects kustomize navigates,
// indexes them by group-version-kind, and installs the precomputed
// namespace-scope table.
func parseCompactBuiltinSchema(b []byte) error {
	var bs kubernetesapi.BuiltinSchema
	if err := json.Unmarshal(b, &bs); err != nil {
		return errors.Wrap(err)
	}

	defs := spec.Definitions{}
	for name, node := range bs.Definitions {
		defs[name] = rebuildSchema(node)
	}
	AddDefinitions(defs)

	if globalSchema.namespaceabilityByResourceType == nil {
		globalSchema.namespaceabilityByResourceType = make(map[yaml.TypeMeta]bool)
	}
	for _, s := range bs.Scopes {
		apiVersion := s.Version
		if s.Group != "" {
			apiVersion = s.Group + "/" + s.Version
		}
		globalSchema.namespaceabilityByResourceType[yaml.TypeMeta{
			APIVersion: apiVersion,
			Kind:       s.Kind,
		}] = s.Namespaced
	}
	return nil
}

// rebuildSchema reconstructs a kube-openapi spec.Schema from a compact node,
// reproducing exactly the fields and extension value types kustomize reads
// (the x-kubernetes-* extensions are rebuilt with the dynamic types the
// ResourceSchema accessors assert on: string, []interface{}, and
// []interface{} of map[string]interface{}).
func rebuildSchema(n *kubernetesapi.SchemaNode) spec.Schema {
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
		ext[kubernetesPatchStrategyExtensionKey] = n.PatchStrategy
	}
	if n.PatchMergeKey != "" {
		ext[kubernetesMergeKeyExtensionKey] = n.PatchMergeKey
	}
	if len(n.ListMapKeys) > 0 {
		keys := make([]interface{}, len(n.ListMapKeys))
		for i, k := range n.ListMapKeys {
			keys[i] = k
		}
		ext[kubernetesMergeKeyMapList] = keys
	}
	if len(n.GroupVersionKinds) > 0 {
		gvks := make([]interface{}, len(n.GroupVersionKinds))
		for i, g := range n.GroupVersionKinds {
			gvks[i] = map[string]interface{}{
				groupKey:   g.Group,
				versionKey: g.Version,
				kindKey:    g.Kind,
			}
		}
		ext[kubernetesGVKExtensionKey] = gvks
	}
	if len(ext) > 0 {
		s.Extensions = ext
	}
	return s
}
