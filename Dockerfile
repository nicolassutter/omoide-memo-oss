# syntax=docker/dockerfile:1.7

ARG GO_VERSION
ARG NODE_VERSION
ARG PNPM_VERSION

FROM golang:${GO_VERSION}-alpine AS api-builder
WORKDIR /src/packages/telemetry-api
COPY packages/telemetry-api/go.mod packages/telemetry-api/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY packages/telemetry-api/ ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go run . openapi > /tmp/openapi.yaml
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

FROM node:${NODE_VERSION}-alpine AS dashboard-builder
RUN npm install -g "pnpm@${PNPM_VERSION}"
WORKDIR /src
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages/ ./packages/
COPY --from=api-builder /tmp/openapi.yaml /src/packages/telemetry-api/doc/openapi.yaml
WORKDIR /src/packages/telemetry-dashboard
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile
RUN pnpm --filter telemetry-dashboard generate:openapi
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm --filter telemetry-dashboard build

FROM node:${NODE_VERSION}-alpine AS dashboard-runtime
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0
COPY --from=dashboard-builder --chown=node:node /src/packages/telemetry-dashboard/.output/ ./
USER node
EXPOSE 3000
CMD ["node", "server/index.mjs"]
