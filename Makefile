

WORKDIR ?= cmd/ingress-proxy

up: tidy
	cd $(WORKDIR) && go run .

tidy:
	go mod tidy