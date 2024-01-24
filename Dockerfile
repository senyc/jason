FROM docker.io/golang:1.21 AS builder

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
WORKDIR /build/cmd/jason
RUN CGO_ENABLED=0 go build -o jason-run

FROM docker.io/alpine:3
COPY --from=builder /build/cmd/jason/jason-run /jason-run

EXPOSE 8080

ENTRYPOINT ["/jason-run"]
