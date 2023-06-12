.DEFAULT_GOAL := help
target=dev
tag=latest
service=client

## —— 📦 Makefile 📦 —————————————————————————————————————————

help: ## Outputs this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' Makefile | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

## —— 🐋 Docker 🐋 —————————————————————————————————————————

build: ## Up the docker environment
	target=${target} docker-compose build --progress=plain
up: ## Up the docker environment
	target=${target} docker-compose up -d
down: ## Down the docker environment
	target=${target} docker-compose down
release:
	target=prod docker-compose build --no-cache --progress=plain
run-prod:
	docker run --rm --name app-prod go-dev:19.2-prod

## —— Go Container —————————————————————————————————————————
enter: ## Access to manager container
	docker exec -u docker-user -ti pfm-${service} /usr/bin/fish
run: ## run main.go
	docker exec -ti go-dev go run main.go
debug: ## run on debug mode, remember to listen on the IDE
	docker exec -ti go-manager dlv --headless --listen=:40000 --api-version=2 debug
debug-test: ## run on debug mode, remember to listen on the IDE
	docker exec -ti go-manager dlv --headless --listen=:40000 --api-version=2 test
test: ##@tests Start all or <container_name> containers in foreground
	cd tests; go test ./... -v; cd -
coverage:
	go test -covermode=atomic -coverprofile=cover.out ./... -v -coverpkg=./... && go tool cover -html=cover.out -o cover.html

## —— Proto —————————————————————————————————————————
protoGen: ## @proto Start all or <container_name> containers in foreground
	protoc 	app/Adapter/In/ApiGrcp/proto/*.proto --go_out=app/Adapter/In/ApiGrcp/gen --go-grpc_out=app/Adapter/In/ApiGrcp/gen

## —— Evans —————————————————————————————————————————
evans: ## evans
	./evans.sh -p 9090 -r

## —— 🐋 Git 🐋 —————————————————————————————————————————
clearGit:

	git branch --merged | grep -v "master" >/tmp/merged-branches && vi /tmp/merged-branches && xargs git branch -d </tmp/merged-branches

confirm:
	@( read -p "$(RED)Are you sure? [y/N]$(RESET): " sure && case "$$sure" in [yY]) true;; *) false;; esac )

