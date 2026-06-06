// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package openapi_test

import (
	"fmt"

	"sigs.k8s.io/kustomize/kyaml/openapi"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func Example() {
	s := openapi.SchemaForResourceType(yaml.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"})

	f := s.Lookup("spec", "replicas")
	fmt.Println(f.Schema.Type)

	// Output:
	// [integer]
}

func Example_arrayMerge() {
	s := openapi.SchemaForResourceType(yaml.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"})

	f := s.Lookup("spec", "template", "spec", "containers")
	fmt.Println(f.Schema.Type)
	fmt.Println(f.PatchStrategyAndKey()) // merge patch strategy on name

	// Output:
	// [array]
	// merge name
}

func Example_arrayReplace() {
	s := openapi.SchemaForResourceType(yaml.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"})

	f := s.Lookup("spec", "template", "spec", "containers", openapi.Elements, "args")
	fmt.Println(f.Schema.Type)
	ps, mk := f.PatchStrategyAndKey() // no patch strategy or merge key
	fmt.Printf("strategy=%q key=%q\n", ps, mk)

	// Output:
	// [array]
	// strategy="" key=""
}

func Example_arrayElement() {
	s := openapi.SchemaForResourceType(yaml.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"})

	f := s.Lookup("spec", "template", "spec", "containers",
		openapi.Elements, "ports", openapi.Elements, "containerPort")
	fmt.Println(f.Schema.Type)

	// Output:
	// [integer]
}

func Example_map() {
	s := openapi.SchemaForResourceType(yaml.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"})

	f := s.Lookup("metadata", "labels")
	fmt.Println(f.Schema.Type)

	// Output:
	// [object]
}
