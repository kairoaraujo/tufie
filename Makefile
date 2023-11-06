.PHONY: test
test:
	go test -race -covermode atomic -coverprofile=cover.out ./... -v

.PHONY: coverage
coverage:
	go tool cover -html=cover.out


.PHONY: lint
coverage:
	golangci-lint run
	go vet
