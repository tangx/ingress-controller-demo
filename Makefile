

WORKDIR ?= cmd/squid

up: tidy
	cd $(WORKDIR) && env=local go run .

tidy:
	go mod tidy