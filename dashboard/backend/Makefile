GO_TEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

go-mod:
	go mod download

go-test:
	go test -race -v -cover -coverprofile coverage.out $(GO_TEST_PACKAGES)

go-lint:
	golangci-lint run -v

