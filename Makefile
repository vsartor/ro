
install:
	go install -ldflags "-X main.Version=$(shell git rev-list --all --count)"

tag:
	git tag -a "1r$(shell git rev-list --all --count)"
