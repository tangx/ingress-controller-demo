FROM golang:1.17 as env
ENV GOPROXY=https://goproxy.cn,direct

FROM env as builder
WORKDIR /src
ADD . .
RUN cd cmd/squid/ && CGO_ENABLED=0 go build -o /go/bin/squid .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/bin/squid /go/bin/squid
ENTRYPOINT [ "/go/bin/squid" ]
