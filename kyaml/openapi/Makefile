# Copyright 2020 The Kubernetes Authors.
# SPDX-License-Identifier: Apache-2.0

MYGOBIN = $(shell go env GOBIN)
ifeq ($(MYGOBIN),)
MYGOBIN = $(shell go env GOPATH)/bin
endif
KIND_VERSION := "v0.11.1"
API_VERSION ?= "v1.21.2"

.PHONY: all
all: \
	kubernetesapi/schema.json.gz \
	kubernetesapi/openapiinfo.go

.PHONY: clean
clean:
	rm -f kubernetesapi/openapiinfo.go

# This will remove all currently built-in schema,
# so think twice before deleting.
# To replace what this will delete typically requires the ability
# to contact a live kubernetes API server.
.PHONY: nuke
nuke: clean
	rm -rf kubernetesapi/v1_*

$(MYGOBIN)/kind:
	( \
		set -e; \
		d=$(shell mktemp -d); cd $$d; \
		wget -O ./kind https://github.com/kubernetes-sigs/kind/releases/download/$(KIND_VERSION)/kind-$(shell uname)-amd64; \
		chmod +x ./kind; \
		mv ./kind $(MYGOBIN); \
		rm -rf $$d; \
	)

.PHONY: kubernetesapi/openapiinfo.go
kubernetesapi/openapiinfo.go:
	./scripts/makeOpenApiInfoDotGo.sh

.PHONY: kubernetesapi/swagger.pb
kubernetesapi/swagger.pb: $(MYGOBIN)/kind $(MYGOBIN)/kustomize
	./scripts/fetchSchemaFromCluster.sh $(API_VERSION)

# schema.json.gz is the embedded, committed built-in schema: the compact set of
# fields kustomize reads, extracted from swagger.pb (see kubernetesapi/scripts/
# extractschema) and gzip-compressed. It is regenerated from a fetched swagger.pb.
.PHONY: kubernetesapi/schema.json.gz
kubernetesapi/schema.json.gz: kubernetesapi/swagger.pb
	./scripts/extractSchema.sh $(API_VERSION)

$(MYGOBIN)/kustomize:
	$(shell cd ../.. && MYGOBIN=$(MYGOBIN) make $(MYGOBIN)/kustomize)
