GO_PATH=$(shell go env GOPATH)

install_deps:
	@echo "Installing go deps"
	@go mod tidy
	@echo "Installing linters..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_PATH)/bin v1.53.3
	go install golang.org/x/tools/cmd/goimports@latest

	@echo "Installing report tools"
	go install github.com/mcubik/goverreport@latest

	@echo "Installing pre-commit tools"
	pip3 install pre-commit
	@echo If "(pre-commit install)" fails try to restar your terminal and run it manually
	pre-commit install

lint:
	@echo "Running linters"
	gofmt -w . && goimports -w .
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 --config=.golangci.yaml --new-from-rev=HEAD~1 --verbose ./pkg/...

pre_coverage_tests:
	@./scripts/test-db-up.sh

coverage_tests: pre_coverage_tests
	@echo "Running tests"
	go clean -testcache
	go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1
	@make post_coverage_tests
	goverreport -coverprofile=/tmp/coverage.out -sort=block -order=desc -threshold=90 || (echo -e "**********Minimum test coverage was not reached(90%)**********"; exit 1)
	go tool cover -html=/tmp/coverage.out

post_coverage_tests:
	@./scripts/test-db-down.sh