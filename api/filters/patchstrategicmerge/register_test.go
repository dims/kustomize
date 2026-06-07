// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package patchstrategicmerge_test

// These tests exercise schema-driven strategic merge, which is now opt-in;
// register the built-in Kubernetes schema for the test binary.
import _ "sigs.k8s.io/kustomize/kyaml/openapi/builtinschema"
