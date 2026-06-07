// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package krusty_test

// These tests exercise the built-in Kubernetes schema (strategic merge, namespace
// scoping), which is now opt-in; register it for the test binary.
import _ "sigs.k8s.io/kustomize/kyaml/openapi/builtinschema"
