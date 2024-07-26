PREFIX?=$(shell pwd)

NAME := gomod
BINDIR := ${PREFIX}/bin

PKG := github.com/cfanbo/gomod

GO111MODULE=on
CGO_ENABLED := 0

# Set any default go build tags
BUILDTAGS :=
BUILDMETA :=

# Populate version variables, Add to compile time flags
VERSION := $(shell git describe --tags `git rev-list --tags --max-count=1`)
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	BUILDMETA := dirty
endif

CREATORVAR=-X ${PKG}/core.GitCommit=$(GITCOMMIT) \
	-X ${PKG}/core.Version=$(VERSION) \
	-X ${PKG}/core.BuildMeta=$(BUILDMETA) \

GO ?= "go"
GO_LDFLAGS=-ldflags "-s -w $(CREATORVAR)"

.PHONY: all test clean build install

all: test install

build:
	@echo "==> $@"
	@CGO_ENABLED=0 GO111MODULE=${GO111MODULE} $(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(BINDIR)/$(NAME) main.go

install:
	@echo "==> $@"
	@go install

test:
	go test ./...

clean:
	$(RM) -r $(BINDIR)
	go clean -i ./...