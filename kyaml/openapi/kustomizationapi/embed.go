// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package kustomizationapi

import _ "embed"

//go:embed swagger.json
var swaggerJSON []byte

// MustAsset returns the built-in Kustomization OpenAPI schema as JSON. The name
// argument is ignored, as there is a single embedded asset.
func MustAsset(string) []byte {
	return swaggerJSON
}
