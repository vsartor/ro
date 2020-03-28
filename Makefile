
install:
	go install -ldflags "-X main.Version=$(shell git rev-list --all --count)"
