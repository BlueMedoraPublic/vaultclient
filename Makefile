ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

export LOCAL_VAULT_CONTAINER_NAME := dev-vault
export LOCAL_VAULT_TOKEN := token
export LOCAL_VAULT_PORT := 38200
export LOCAL_VAULT_ADDR := http://localhost:38200

all: fmt lint test

test: local-vault
	go test ./...

lint:
	golint ./...

fmt:
	go fmt ./...

local-vault: clean-vault
	@scripts/local_vault.sh

clean-vault:
	$(shell docker rm dev-vault --force >> /dev/null 2>&1 || true)
