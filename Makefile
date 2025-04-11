GHCR_REPO:=ghcr.io/babbage88/smbplusplus:
GHCR_REPO_TEST:=jtrahan88/smbp2-test:
DEPLOYMENT:=deployment/k8s/smbplus2.yaml
export BUILDER:=smbplusplus-builder
ENV_FILE:=.env
MIG:=$(shell date '+%m%d%Y.%H%M%S')
SHELL := /bin/bash
export goose_src:=internal/goose/goose.go
export goose_out_bin:=goosey
export svc:=smbplus2
export pgpw:=none
export pguser:=jtrahan
export dbname:=smbplusplus
export schema_dump_file:=schema.sql

check-builder:
	@if ! docker buildx inspect $(BUILDER) > /dev/null 2>&1; then \
		echo "Builder $(BUILDER) does not exist. Creating..."; \
    	docker buildx create --name $(BUILDER) --bootstrap; \
	fi

create-builder: check-builder

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models && swagger generate spec -o swagger.json --scan-models

dump-schema:
	PGPASSWORD=$(pgpw) pg_dump -h localhost -p 5432 -s -U $(pguser) $(dbname) > $(schema_dump_file)

build-goose:
	$(info ************  BUILD GOOSE MIGRATION BINARY src: $(goose_src) dest: $(goose_out_bin)  ************)
	go build -v -o $(goose_out_bin) $(goose_src)

dev-swagger: check-swagger
	$(info ************ GENERATING SWAGGER SPEC: dev development ************)
	swagger generate spec -o ./dev-swagger.yaml --scan-models && swagger generate spec -o dev-swagger.json --scan-models
	swagger mixin spec/swagger.dev.json dev-swagger.json --output swagger.json --format=json
	swagger mixin spec/swagger.dev.yaml dev-swagger.yaml --output swagger.yaml --format=yaml
	rm dev-swagger.json && rm dev-swagger.yaml

local-swagger: check-swagger
	$(info ************ GENERATING SWAGGER SPEC: local development ************)
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
	$(info ************ Starting application on localshost: go run . ************)
	go run . --development

embed-swagger:
	swagger generate spec -o ./embed/swagger.yaml --scan-models && swagger generate spec > ./embed/swagger.json

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml --no-open --port 4443

buildandpushdev: dev-swagger
	$(info ************ Performing docker buildx buildandpush Repository: $(GHCR_REPO)$(tag) ************)
	docker buildx use $(BUILDER)
	docker buildx build --platform linux/amd64,linux/arm64 -t $(GHCR_REPO)$(tag) . --push

buildandpushlocalk3: k3local-swagger
	docker buildx use $(BUILDER)
	docker buildx build --platform linux/amd64,linux/arm64 -t $(GHCR_REPO_TEST)$(tag) . --push

deploydev: buildandpushdev
	$(info ************ Applying kubernetes manifest $(DEPLOYMENT) and restarting service: $(svc) ************)
	kubectl apply -f $(DEPLOYMENT)
	kubectl rollout restart deployment $(svc)

deploylocalk3: buildandpushlocalk3
	$(info ************ Applying kubernetes manifest $(LOCALK3DEPLOYMENT) and restarting service: $(svc) ************)
	kubelocal apply -f $(LOCALK3DEPLOYMENT)
	kubelocal rollout restart deployment $(svc)

new-sqlmigration:
	@source $(ENV_FILE) && \
		export GOOSE_DBSTRING GOOSE_MIGRATION_DIR GOOSE_DRIVER && \
		goose create -s $(MIG) sql

