GHCR_REPO:=ghcr.io/babbage88/smbp2:
GHCR_REPO_TEST:=jtrahan88/smbp2-test:
DEPLOYMENT:=deployment/k8s/smbplus2.yaml
BUILDER:=smbplusplus-builder
ENV_FILE:=.env
MIG:=$(shell date '+%m%d%Y.%H%M%S')
SHELL := /bin/bash
export svc:=smbplus2

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models && swagger generate spec -o swagger.json --scan-models

dev-swagger: check-swagger
	swagger generate spec -o ./dev-swagger.yaml --scan-models && swagger generate spec -o dev-swagger.json --scan-models
	swagger mixin spec/swagger.dev.json dev-swagger.json --output swagger.json --format=json
	swagger mixin spec/swagger.dev.yaml dev-swagger.yaml --output swagger.yaml --format=yaml
	rm dev-swagger.json && rm dev-swagger.yaml

local-swagger: check-swagger
	swagger generate spec -o ./local-swagger.yaml --scan-models && swagger generate spec --scan-models -o local-swagger.json --scan-models
	swagger mixin spec/swagger.local.json local-swagger.json --output swagger.json --format=json
	swagger mixin spec/swagger.local.yaml local-swagger.yaml --output swagger.yaml --format=yaml
	rm local-swagger.json && rm local-swagger.yaml

k3local-swagger: check-swagger
	swagger generate spec -o ./k3local-swagger.yaml --scan-models && swagger generate spec -o k3local-swagger.json --scan-models
	swagger mixin spec/swagger.localdev.json k3local-swagger.json --output swagger.json --format=json
	swagger mixin spec/swagger.localdev.yaml k3local-swagger.yaml --output swagger.yaml --format=yaml
	rm k3local-swagger.json && rm k3local-swagger.yaml

run-local: local-swagger
	go run .

embed-swagger:
	swagger generate spec -o ./embed/swagger.yaml --scan-models && swagger generate spec > ./embed/swagger.json

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml --no-open --port 4443

buildandpushdev: dev-swagger
	docker buildx use $(BUILDER)
	docker buildx build --platform linux/amd64,linux/arm64 -t $(GHCR_REPO)$(tag) . --push

buildandpushlocalk3: k3local-swagger
	docker buildx use $(BUILDER)
	docker buildx build --platform linux/amd64,linux/arm64 -t $(GHCR_REPO_TEST)$(tag) . --push

deploydev: buildandpushdev
	kubectl apply -f $(DEPLOYMENT)
	kubectl rollout restart deployment $(svc)

deploylocalk3: buildandpushlocalk3
	kubelocal apply -f $(LOCALK3DEPLOYMENT)
	kubelocal rollout restart deployment $(svc)

new-sqlmigration:
	@source $(ENV_FILE) && \
		export GOOSE_DBSTRING GOOSE_MIGRATION_DIR GOOSE_DRIVER && \
		goose create -s $(MIG) sql

