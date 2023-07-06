FROM golang:1.19-alpine AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /work
COPY . /work

# Build admission-webhook
RUN go build -o bin/go-webhook cmd/api/*

# ---
FROM ubuntu:23.04 AS run
ARG VERSION
ENV VERSION=${VERSION}


COPY --from=build /work/bin/go-webhook /usr/local/bin/
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY ./email-templates  /work/email-templates

CMD ["app"]
