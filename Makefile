

WORKDIR ?= cmd/squid

up: tidy
	cd $(WORKDIR) && go run .

tidy:
	go mod tidy