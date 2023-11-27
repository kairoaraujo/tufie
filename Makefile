.PHONY: test
test:
	go test -covermode atomic -coverprofile=cover.out ./... -v

.PHONY: coverage
coverage:
	go tool cover -html=cover.out


.PHONY: lint
lint:
	golangci-lint run
	go vet
