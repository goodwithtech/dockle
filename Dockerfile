FROM golang:1.16-alpine AS builder
WORKDIR /go/src/github.com/Portshift/dockle/
COPY ./ ./
RUN CGO_ENABLED=0 go build -o dockle_remote ./cmd/dockle_remote/

FROM registry.access.redhat.com/ubi8
RUN yum install ca-certificates -y
RUN mkdir /licenses
COPY ./LICENSE /licenses/
RUN mkdir /app
COPY --from=builder /go/src/github.com/Portshift/dockle/dockle_remote /app/dockle
RUN chmod +x /app/dockle
ENTRYPOINT ["/app/dockle"]

# Build-time metadata as defined at http://label-schema.org
ARG BUILD_DATE
ARG VCS_REF
LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="dockle" \
    org.label-schema.description="CIS Docker benchmark for images stored in a private or public Docker registries using Dockle" \
    org.label-schema.url="https://github.com/Portshift/dockle" \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vcs-url="https://github.com/Portshift/dockle"

### Required OpenShift Labels
ARG IMAGE_VERSION
LABEL name="dockle" \
      vendor="Portshift" \
      version=${IMAGE_VERSION} \
      release=${IMAGE_VERSION} \
      summary="CIS Docker benchmark" \
      description="CIS Docker benchmark for images stored in a private or public Docker registries using Dockle"

USER 1000