// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package yaml_test

// These tests exercise the built-in Kubernetes schema (via the openapi package),
// which is now opt-in; register it for the test binary. This is an external test
// package because the openapi package imports kyaml/yaml.
import _ "sigs.k8s.io/kustomize/kyaml/openapi/builtinschema"
