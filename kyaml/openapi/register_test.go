// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package openapi

import (
	"os"
	"testing"

	"sigs.k8s.io/kustomize/kyaml/openapi/kubernetesapi"
)

// TestMain registers the built-in Kubernetes schema for this package's tests.
// The openapi package no longer imports the schema data (it is opt-in); tests
// that exercise the built-in schema register it here. Importing the data package
// from this test is cycle-free because kubernetesapi does not import openapi.
func TestMain(m *testing.M) {
	RegisterBuiltinSchema(kubernetesapi.Schema())
	os.Exit(m.Run())
}
