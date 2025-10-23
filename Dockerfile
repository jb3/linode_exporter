ARG golang_version=1.25.3
ARG os_version=alpine
FROM golang:$golang_version-$os_version AS builder
    LABEL org.opencontainers.image.authors="Joe Banks <joe@jb3.dev>"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/linode_exporter.go

FROM $os_version:latest

WORKDIR /

COPY --from=builder /app/linode_exporter /linode_exporter

ENTRYPOINT ["/linode_exporter"]
