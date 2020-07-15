FROM golang:alpine

RUN apk add make

WORKDIR /stories

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

EXPOSE 8080

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

CMD ["./out/stories", "serve"]