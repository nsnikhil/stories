FROM golang:alpine as builder
WORKDIR /stories
ARG SSH_PRIVATE_KEY
RUN apk update &&\
    apk add --no-cache openssh-client &&\
    apk add git &&\
    mkdir -p /root/.ssh/ &&\
    echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa &&\
    chmod 600 /root/.ssh/id_rsa &&\
    touch /root/.ssh/known_hosts &&\
    ssh-keyscan -t rsa github.com >> /root/.ssh/known_hosts &&\
    git config --global --add url."git@github.com:".insteadOf "https://github.com/"
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o stories cmd/*.go

FROM scratch
COPY --from=builder /stories/stories .
CMD ["./stories", "grpc-serve"]