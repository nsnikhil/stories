APP=stories
APP_VERSION:=0.2
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

setup: copy-config init-db migrate test

init-db:
	psql -c "create user stories_user superuser password 'stories_password';" -U postgres
	psql -c "create database stories_db owner=stories_user" -U postgres

deps:
	go mod download

tidy:
	go mod tidy

check: fmt vet lint

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	golint $(ALL_PACKAGES)

compile:
	mkdir -p out/
	go build -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" -o $(APP_EXECUTABLE) cmd/*.go

build: deps compile

docker-build:
	docker build --build-arg SSH_PRIVATE_KEY="$$(cat ~/.ssh/travis_ci_key)" -t nsnikhil/$(APP):$(APP_VERSION) .
	docker rmi -f $$(docker images -f "dangling=true" -q)

ci-docker-build:
	docker build --build-arg SSH_PRIVATE_KEY="$$(cat ~/.ssh/id_rsa)" -t nsnikhil/$(APP):$(APP_VERSION) .
	docker rmi -f $$(docker images -f "dangling=true" -q)

docker-push: docker-build
	docker push nsnikhil/$(APP):$(APP_VERSION)

ci-docker-push: ci-docker-build
	docker push nsnikhil/$(APP):$(APP_VERSION)

serve: build
	$(APP_EXECUTABLE) serve

clean:
	rm -rf out/

copy-config:
	cp env.sample env.yaml

test:
	go clean -testcache
	go test ./...

ci-test: copy-config init-db migrate test

test-cover-html:
	go clean -testcache
	mkdir -p out/
	go test ./... -coverprofile=out/coverage.out
	go tool cover -html=out/coverage.out

migrate: build
	$(APP_EXECUTABLE) migrate

rollback: build
	$(APP_EXECUTABLE) rollback