FROM golang:alpine as builder
WORKDIR /stories
RUN apk update && apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o stories cmd/*.go

FROM scratch
COPY --from=builder /stories/stories .
CMD ["./stories", "grpc-serve"]