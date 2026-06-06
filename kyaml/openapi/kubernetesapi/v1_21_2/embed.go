// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package v1_21_2

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"sync"
)

//go:embed schema.json.gz
var schemaJSONGz []byte

var (
	once       sync.Once
	schemaJSON []byte
)

// MustAsset returns the compact built-in Kubernetes schema for this version as
// JSON (see kubernetesapi.BuiltinSchema), decompressed once from the embedded
// gzip. The name argument is ignored, as there is a single embedded asset.
func MustAsset(string) []byte {
	once.Do(func() {
		r, err := gzip.NewReader(bytes.NewReader(schemaJSONGz))
		if err != nil {
			panic(err)
		}
		defer r.Close()
		if schemaJSON, err = io.ReadAll(r); err != nil {
			panic(err)
		}
	})
	return schemaJSON
}
