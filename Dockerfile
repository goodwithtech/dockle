FROM golang:1.22-alpine AS builder
ARG TARGET_DIR="github.com/goodwithtech/dockle"

WORKDIR /go/src/${TARGET_DIR}
RUN apk --no-cache add git
COPY . .
RUN CGO_ENABLED=0 go build -a -o /dockle ${PWD}/cmd/dockle

FROM alpine:3.17
COPY --from=builder /dockle /usr/local/bin/dockle
RUN chmod +x /usr/local/bin/dockle
RUN apk --no-cache add ca-certificates shadow

# for use docker daemon via mounted /var/run/docker.sock
#RUN addgroup -S docker && adduser -S -G docker dockle && usermod -aG root dockle
#USER dockle

ENTRYPOINT ["dockle"]
