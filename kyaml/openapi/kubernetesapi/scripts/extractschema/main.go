// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Command extractschema reads a full Kubernetes OpenAPI v2 protobuf (swagger.pb)
// and writes the compact, purpose-built schema kustomize needs
// (kubernetesapi.BuiltinSchema) as gzip-compressed JSON.
//
// It parses the proto with the same gnostic + kube-openapi path the openapi
// package uses (proto.Unmarshal -> openapi_v2.Document -> spec.Swagger), so the
// extracted values match what kustomize would have seen, then keeps only:
//   - per definition: type, $ref, items, properties, additionalProperties, and
//     the x-kubernetes-{patch-strategy,patch-merge-key,list-map-keys,
//     group-version-kind} extensions;
//   - per resource type: whether it is namespace-scoped (the same determination
//     openapi.findNamespaceability makes from the paths).
package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	openapi_v2 "github.com/google/gnostic-models/openapiv2"
	"google.golang.org/protobuf/proto"
	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/kustomize/kyaml/openapi/kubernetesapi"
)

const (
	gvkKey = "x-kubernetes-group-version-kind"
	psKey  = "x-kubernetes-patch-strategy"
	mkKey  = "x-kubernetes-patch-merge-key"
	lmkKey = "x-kubernetes-list-map-keys"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: extractschema <input swagger.pb> <output schema.json.gz>")
		os.Exit(2)
	}

	raw, err := os.ReadFile(os.Args[1])
	check(err)
	doc := &openapi_v2.Document{}
	check(proto.Unmarshal(raw, doc))
	var sw spec.Swagger
	_, err = sw.FromGnostic(doc)
	check(err)

	out := &kubernetesapi.BuiltinSchema{Definitions: map[string]*kubernetesapi.SchemaNode{}}
	for name, def := range sw.Definitions {
		out.Definitions[name] = projectNode(def)
	}
	out.Scopes = extractScopes(&sw)

	b, err := json.Marshal(out)
	check(err)

	f, err := os.Create(os.Args[2])
	check(err)
	zw, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	check(err)
	_, err = zw.Write(b)
	check(err)
	check(zw.Close())
	check(f.Close())

	st, _ := os.Stat(os.Args[2])
	fmt.Printf("definitions=%d scopes=%d json=%d bytes gz=%d bytes\n",
		len(out.Definitions), len(out.Scopes), len(b), st.Size())
}

// projectNode keeps only the schema fields kustomize navigates or merges against.
func projectNode(s spec.Schema) *kubernetesapi.SchemaNode {
	n := &kubernetesapi.SchemaNode{}
	if len(s.Type) > 0 {
		n.Type = append([]string(nil), s.Type...)
	}
	if r := s.Ref.String(); r != "" {
		n.Ref = strings.TrimPrefix(r, "#/definitions/")
	}
	if s.Items != nil && s.Items.Schema != nil {
		n.Items = projectNode(*s.Items.Schema)
	}
	if len(s.Properties) > 0 {
		n.Properties = make(map[string]*kubernetesapi.SchemaNode, len(s.Properties))
		for k, v := range s.Properties {
			n.Properties[k] = projectNode(v)
		}
	}
	if s.AdditionalProperties != nil && s.AdditionalProperties.Schema != nil {
		n.AdditionalProperties = projectNode(*s.AdditionalProperties.Schema)
	}
	if v, ok := s.Extensions[psKey].(string); ok {
		n.PatchStrategy = v
	}
	if v, ok := s.Extensions[mkKey].(string); ok {
		n.PatchMergeKey = v
	}
	if v, ok := s.Extensions[lmkKey]; ok {
		n.ListMapKeys = toStringSlice(v)
	}
	if v, ok := s.Extensions[gvkKey]; ok {
		n.GroupVersionKinds = toGVKs(v)
	}
	return n
}

// extractScopes mirrors openapi.findNamespaceability: a resource type is
// namespace-scoped if any of its GET paths contains a namespace path parameter.
func extractScopes(sw *spec.Swagger) []kubernetesapi.ScopeEntry {
	type key struct{ g, v, k string }
	ns := map[key]bool{}
	if sw.Paths != nil {
		for path, pi := range sw.Paths.Paths {
			if pi.Get == nil {
				continue
			}
			ext, ok := pi.Get.Extensions[gvkKey].(map[string]interface{})
			if !ok {
				continue
			}
			kk := key{getStr(ext, "group"), getStr(ext, "version"), getStr(ext, "kind")}
			if strings.Contains(path, "namespaces/{namespace}") {
				ns[kk] = true
			} else if _, exists := ns[kk]; !exists {
				ns[kk] = false
			}
		}
	}
	out := make([]kubernetesapi.ScopeEntry, 0, len(ns))
	for k, v := range ns {
		out = append(out, kubernetesapi.ScopeEntry{Group: k.g, Version: k.v, Kind: k.k, Namespaced: v})
	}
	return out
}

func toStringSlice(v interface{}) []string {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, e := range list {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func toGVKs(v interface{}) []kubernetesapi.GVK {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]kubernetesapi.GVK, 0, len(list))
	for _, e := range list {
		m, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, kubernetesapi.GVK{
			Group:   getStr(m, "group"),
			Version: getStr(m, "version"),
			Kind:    getStr(m, "kind"),
		})
	}
	return out
}

func getStr(m map[string]interface{}, k string) string {
	if s, ok := m[k].(string); ok {
		return s
	}
	return ""
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "extractschema:", err)
		os.Exit(1)
	}
}
