VERSION      := $(shell git describe --always --tags --abbrev=6 --dirty="-dev")
VERSION_FLAG := -X "main.Version=$(VERSION)"
LD_FLAGS     := -ldflags='$(VERSION_FLAG)' # -linkmode external -extldflags "-static"'
MOD_FLAG     := -mod=vendor

.PHONY: all
all: ogrt

.PHONY: ogrt
ogrt:
	$Q go build $(MOD_FLAG) $(LD_FLAGS) -o bin/ogrt cmd/ogrt/main.go

test:
	$Q go test $(MOD_FLAG) ./...

.PHONY: fmt
fmt:
	$Q go fmt ./...

Q := $(if $V,,@)
