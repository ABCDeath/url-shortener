include .env

PROJECTNAME := $(shell basename "$(PWD)")

PROJECT_DIR := $(shell pwd)

BACKEND_SRC_DIR := $(PROJECT_DIR)/backend
BACKEND_BUILD_DIR := $(BACKEND_SRC_DIR)/build
BACKEND_BINARY_DEBUG := $(BACKEND_BUILD_DIR)/debug
BACKEND_BINARY_RELEASE := $(BACKEND_BUILD_DIR)/release
BACKEND_SRC := $(shell find $(BACKEND_SRC_DIR)/application -type f -name '*.go' -not -path "./vendor/*")

FRONTEND_SRC := $(PROJECT_DIR)/frontend
FRONTEND_BUILD_DIR := $(FRONTEND_SRC)/dist

# common

.PHONY: all
all: $(BACKEND_BINARY_RELEASE) frontend
	@true

.PHONY: clean
clean:
	@rm $(BACKEND_BUILD_DIR)/*
	@rm -rf $(FRONTEND_SRC)/dist

# backend

$(BACKEND_BINARY_DEBUG): $(BACKEND_SRC)
	@go build -o $(BACKEND_BINARY_DEBUG) $(BACKEND_SRC)

$(BACKEND_BINARY_RELEASE): $(BACKEND_SRC)
	@go build -ldflags "-s -w" -o $(BACKEND_BINARY_RELEASE) $(BACKEND_SRC)

.PHONY: backend
backend: $(BACKEND_BINARY_DEBUG)
	@true

.PHONY: clean_backend
clean-backend:
	@rm $(BACKEND_BUILD_DIR)/*

.PHONY: test
test: backend
	@govendor test +local

.PHONY: run
run: $(BACKEND_BINARY_RELEASE)
	@GIN_MODE=release $(BACKEND_BINARY_DEBUG)

.PHONY: run_debug
run-debug: $(BACKEND_BINARY_DEBUG)
	@$(BACKEND_BINARY_DEBUG)

# frontend

.PHONY: frontend
frontend:
	@npm --prefix $(FRONTEND_SRC) run build

.PHONY: run-frontend
run-frontend:
	@npm --prefix $(FRONTEND_SRC) run serve

# docker

.PHONY: docker
docker:
	@docker-compose -f $(PROJECTNAME).yml build

.PHONY: run_db
run-db:
	@docker-compose -f $(PROJECTNAME).yml up -d mongo

.PHONY: run_docker
run-docker:
	@docker-compose -f $(PROJECTNAME).yml up -d

.PHONY: stop_db
stop-docker:
	@docker-compose -f $(PROJECTNAME).yml down
