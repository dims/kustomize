#!/bin/bash
# Copyright 2026 The Kubernetes Authors.
# SPDX-License-Identifier: Apache-2.0

# Extracts the compact, purpose-built built-in schema kustomize needs from
# kubernetesapi/<version>/swagger.pb into the embedded, committed
# kubernetesapi/<version>/schema.json.gz. See kubernetesapi/scripts/extractschema
# for exactly what is kept. The raw swagger.pb is the regeneration source only and
# is not committed (see kubernetesapi/.gitignore).

set -e

VERSION=$1
DIR="kubernetesapi/${VERSION//./_}"

go run ./kubernetesapi/scripts/extractschema "${DIR}/swagger.pb" "${DIR}/schema.json.gz"
