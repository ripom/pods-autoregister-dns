FROM golang:1.11.5-alpine3.8 AS builder

RUN apk add --update --no-cache git ca-certificates

WORKDIR /code

COPY ./dns ./dns/
ENV GOPATH=/code/dns
ENV CGO_ENABLED=0

RUN go build -o /code/dnsprovider /code/dns/main.go

FROM scratch

COPY --from=builder /code/dnsprovider /code/dnsprovider
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/code/dnsprovider"]
