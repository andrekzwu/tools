.PHONY: build clean fmt run check

GOCMD=GO111MODULE=on go

fmt:
	@echo "fmt go files..."
	${GOCMD} fmt ${GOARGS} ./...
	@echo "fmt go files done"
	
stub:
	@echo "stub..."
	go mod tidy -go=1.16 && go mod tidy -go=1.17 && go mod tidy -go=1.18
lint:
	@echo "lint..."
	staticcheck ./...
	golangci-lint run

check:
	@echo "vet code..."
	${GOCMD} vet ./...
	@echo "vet code done"
	@echo "staticcheck code..."
	staticcheck ./...
	@echo "staticcheck code done"