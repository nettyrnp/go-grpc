# Basic go commands
GOCMD      = go
GOBUILD    = $(GOCMD) build
GORUN      = $(GOCMD) run
GOCLEAN    = $(GOCMD) clean
GOTEST     = $(GOCMD) test
GOGET      = $(GOCMD) get
GOFMT      = $(GOCMD) fmt
GOLINT     = $(GOCMD)lint

# GRPC
PROTOC = protoc

#
BINARY1 = ingestor-bi
BINARY2 = persistor-bi
PKGS = ./...

# Texts
TEST_STRING = "TEST"

.PHONY: all help clean test lint format build run grpc

all: clean format lint test build

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binary.'
	@echo '    test               Run all tests.'
	@echo '    lint               Run golint on package sources.'
	@echo '    format             Run gofmt on package sources.'
	@echo '    build              Compile packages and dependencies.'
	@echo '    run                Compile and run Go program.'
	@echo '    grpc               '
	@echo ''
	@echo 'Targets run by default are: clean grpc format lint test.'
	@echo ''

clean:
	@echo "[CLEAN]"
	@$(GOCLEAN)

test:
	@echo "[$(TEST_STRING)]"
	@$(GOTEST) -v $(PKGS)

lint:
	@echo "[LINT]"
	@-$(GOLINT) -min_confidence 0.8 $(PKGS)

format:
	@echo "[FORMAT]"
	@$(GOFMT) $(PKGS)

build:
	@echo "[BUILD]"
	@$(GOBUILD) -o $(BINARY2) ./services/persistor/persistor.go
	@$(GOBUILD) -o $(BINARY1) ./services/ingestor/ingestor.go

run: build
	@echo "[RUN]"
	@$(GORUN) -race ./services/persistor/persistor.go && @$(GORUN) -race ./services/ingestor/ingestor.go

grpc:
	@echo "[GRPC]"
    ${PROTOC} -I ./proto/ ./proto/persistor.proto --go_out=plugins=grpc:./proto


