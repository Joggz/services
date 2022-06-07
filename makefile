SHELL := /bin/bash

# ==============================================================================
# Testing running system

# ./expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
#  hey -m GET -c 100 -n 10000 http://localhost:3000/v1/test
debug: 
	expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
run:
	go run app/services/sales-api/main.go |  go run app/services/tooling/logfmt/main.go

# ==============================================================================
# Building containers

VERSION := 1.0

all: sales-api

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

admin: 
	go run  app/services/tooling/admin/main.go
# ==============================================================================  cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
# Running from within k8s/kind // kubectl config set-context --current --namespace=sales-system
KIND_CLUSTER := joggz-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.23.4@sha256:0e34f0d0fd448aa2f2819cfd74e99fe5793a6e4938b328f657c8e3f81ee0dfb9 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-load:
	cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 |  go run app/services/tooling/logfmt/main.go

kind-update: all kind-load kind-restart

kind-restart:
	kubectl rollout restart deployment sales-pod 

kind-update-apply: all kind-load kind-apply

kind-status-sales:
	kubectl get pods -o wide --watch --all-namespaces

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system
# build:
# 	go build -ldflags "-X main.build=local"

# ==============================================================================
# Modules support
tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck -checks=all ./...
