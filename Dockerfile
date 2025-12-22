ARG GO_VERSION=1.25
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

################################################################################
FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
    ca-certificates \
    tzdata \
    bash \
    postgresql-client \
    && \
    update-ca-certificates

RUN wget -O /usr/local/bin/goose https://github.com/pressly/goose/releases/latest/download/goose_linux_x86_64 \
    && chmod +x /usr/local/bin/goose

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

COPY --from=build /bin/server /bin/
COPY sql/schema/ /app/migrations/
COPY scripts/migrate.sh /scripts/
RUN chmod +x /scripts/migrate.sh

USER appuser
EXPOSE 8080
ENTRYPOINT ["/scripts/migrate.sh"]
