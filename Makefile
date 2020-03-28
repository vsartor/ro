
VERSION := $(shell git rev-list --all --count)

install:
	go install -ldflags "-X main.Version=$(VERSION)"

tag:
	git tag -a "1r$(VERSION)" -m "version $(VERSION)"
