# syntax=docker/dockerfile:1.7

FROM alpine:3.20 AS mise-base
RUN apk add --no-cache curl bash ca-certificates git
ENV MISE_INSTALL_PATH=/usr/local/bin/mise
ENV MISE_DATA_DIR=/mise
ENV MISE_CONFIG_DIR=/mise
ENV MISE_YES=1
RUN curl -fsSL https://mise.run | sh
ENV PATH="/mise/shims:${PATH}"

FROM mise-base AS api-builder
WORKDIR /src/packages/telemetry-api
COPY packages/telemetry-api/mise.toml packages/telemetry-api/go.mod packages/telemetry-api/go.sum ./
RUN mise trust && mise install
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY packages/telemetry-api/ ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go run . openapi > /tmp/openapi.yaml
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api .

FROM alpine:3.20 AS api-runtime
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -g 65532 -S app && adduser -u 65532 -S app -G app
COPY --from=api-builder /out/api /app/api
USER 65532
WORKDIR /app
EXPOSE 9999
ENTRYPOINT ["/app/api"]

FROM mise-base AS dashboard-builder
WORKDIR /src
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages/ ./packages/
COPY --from=api-builder /tmp/openapi.yaml /src/packages/telemetry-api/doc/openapi.yaml
WORKDIR /src/packages/telemetry-dashboard
RUN mise trust && mise install
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile
RUN pnpm --filter telemetry-dashboard generate:openapi
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm --filter telemetry-dashboard build

FROM mise-base AS dashboard-runtime
RUN addgroup -g 1000 -S node && adduser -u 1000 -S node -G node
WORKDIR /app
COPY packages/telemetry-dashboard/mise.toml ./mise.toml
RUN mise trust && mise install node
ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0
COPY --from=dashboard-builder --chown=node:node /src/packages/telemetry-dashboard/.output/ ./
USER node
EXPOSE 3000
CMD ["node", "server/index.mjs"]
