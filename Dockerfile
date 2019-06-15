FROM golang:1.12-alpine AS builder
COPY go.mod go.sum /app/
WORKDIR /app/
RUN apk --no-cache add git
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /dockle cmd/dockle/main.go

FROM alpine:3.9
COPY --from=builder /dockle /usr/local/bin/dockle
RUN chmod +x /usr/local/bin/dockle
RUN apk --no-cache add ca-certificates shadow

# for use docker daemon via mounted /var/run/docker.sock
RUN addgroup -S docker && adduser -S -G docker dockle && usermod -aG root dockle
USER dockle

ENTRYPOINT ["dockle"]
