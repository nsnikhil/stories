APP=stories
APP_VERSION:="0.1"
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

deps:
	go mod download

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

compile:
	mkdir -p out/
	go build -o $(APP_EXECUTABLE) -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" cmd/*.go

build: deps compile

serve: build
	$(APP_EXECUTABLE) serve

k8-serve:
	chmod +x build/kube/start.sh
	./build/kube/start.sh

k8-stop:
	chmod +x build/kube/stop.sh
	./build/kube/stop.sh

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