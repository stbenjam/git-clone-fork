default: fmt vet lint build

build:
	go build -ldflags "${LDFLAGS}" .

install:
	go install -ldflags "${LDFLAGS}"

fmt:
	go fmt .
	git diff --exit-code

lint:
	go run golang.org/x/lint/golint -set_exit_status *.go

vet:
	go vet .

test:
	go test -v ./ironic

clean:
	rm -f terraform-provider-ironic

.PHONY: build install test fmt lint
