// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package merge3_test

// These tests exercise schema-driven strategic merge (via the openapi package),
// which is now opt-in; register the built-in schema for the test binary.
import _ "sigs.k8s.io/kustomize/kyaml/openapi/builtinschema"
