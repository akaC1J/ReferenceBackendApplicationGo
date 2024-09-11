build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -o bin/app ./cmd/app

.PHONY: test-coverage
test-coverage:
	go test -coverprofile=./reports/coverage.out ./... && go tool cover -html=./reports/coverage.out -o ./reports/coverage.html

.PHONY: cyclo cognit lint install-tools

GOPATH := $(shell go env GOPATH)
export PATH := $(GOPATH)/bin:$(PATH)

install-tools:
	@command -v gocyclo >/dev/null 2>&1 || { echo >&2 "Installing gocyclo..."; go install github.com/fzipp/gocyclo/cmd/gocyclo@latest; }
	@command -v gocognit >/dev/null 2>&1 || { echo >&2 "Installing gocognit..."; go install github.com/uudashr/gocognit/cmd/gocognit@latest; }

cyclo: install-tools
	@gocyclo -ignore ".*mock_test.go"  .  > reports/cyclo-report.txt

cognit: install-tools
	@gocognit -ignore ".*mock_test.go" . > reports/cognit-report.txt

lint: cyclo cognit
