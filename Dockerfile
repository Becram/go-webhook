FROM golang:1.19-alpine AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /work
COPY . /work

# Build admission-webhook
RUN go build -o bin/app .

# ---
FROM ubuntu:23.04 AS run
ARG VERSION
ENV ARG VERSION=${VERSION}


COPY --from=build /work/bin/app /usr/local/bin/
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY ./static  /work/static

CMD ["app"]
