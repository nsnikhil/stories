FROM golang:alpine as builder
WORKDIR /stories
RUN apk update && apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o stories cmd/*.go

FROM scratch
COPY --from=builder /stories/stories .
#TODO: REMOVE TEMP HACK OF local.env
COPY local.env .
CMD ["./stories", "grpc-serve", "--configFile", "local.env"]