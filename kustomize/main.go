// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// The kustomize CLI.
package main

import (
	"os"

	"sigs.k8s.io/kustomize/kustomize/v5/commands"

	// Register the built-in Kubernetes OpenAPI schema. The openapi library no
	// longer pulls it in automatically (it is opt-in, so library consumers that
	// don't need it -- e.g. kubectl, which imports only commands/build -- don't
	// vendor the large schema data). The CLI opts in here.
	_ "sigs.k8s.io/kustomize/kyaml/openapi/builtinschema"
)

func main() {
	if err := commands.NewDefaultCommand().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
