FROM golang:1.12-alpine AS builder
ADD go.mod go.sum /app/
WORKDIR /app/
RUN apk --no-cache add git
RUN go mod download
ADD . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /dockle cmd/dockle/main.go

FROM alpine:3.9
COPY --from=builder /dockle /usr/local/bin/dockle
RUN apk --no-cache add ca-certificates
RUN chmod +x /usr/local/bin/dockle

RUN addgroup -S dockle && adduser -S -G dockle dockle
USER dockle
ENTRYPOINT ["dockle"]
