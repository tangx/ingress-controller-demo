

WORKDIR ?= cmd/squid

up: tidy
	cd $(WORKDIR) && env=local go run .

tidy:
	go mod tidy

docker:
	docker build -t cr.docker.tangx.in/controllers/squid:v1.0.01 .
