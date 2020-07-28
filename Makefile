APP=stories
APP_VERSION:=0.1
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

deps:
	go mod download

tidy:
	go mod tidy

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

compile:
	mkdir -p out/
	go build -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" -o $(APP_EXECUTABLE) cmd/*.go

build: deps compile

docker-build:
	docker build -t nsnikhil/$(APP):$(APP_VERSION) .

docker-push: docker-build
	docker push nsnikhil/$(APP):$(APP_VERSION)

serve: build
	$(APP_EXECUTABLE) serve

clean:
	rm -rf out/

test:
	go clean -testcache
	go test ./...

test-cover-html:
	go clean -testcache
	mkdir -p out/
	go test ./... -coverprofile=out/coverage.out
	go tool cover -html=out/coverage.out

migrate: build
	$(APP_EXECUTABLE) migrate

rollback: build
	$(APP_EXECUTABLE) rollback