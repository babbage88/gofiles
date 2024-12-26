# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.0
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# golang dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

# Target go version
ARG TARGETARCH

# Build the application, using cache mount.

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./

# Final stage copy bin and install pre-requisites
FROM alpine:latest AS final

WORKDIR /app

ARG UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  appuser 

RUN mkdir -p /mnt/files/ && \
  chown -R appuser:appuser /app/ && \
  chown -R appuser:appuser /mnt/files
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /app/

# Expose the port that the application listens on.
EXPOSE 8080

ENTRYPOINT [ "/app/server" ]
