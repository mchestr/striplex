# Multi-stage build for both frontend and backend
FROM node:22-slim AS frontend-builder

WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci

# Set production environment for better optimization
ENV NODE_ENV=production \
    GENERATE_SOURCEMAP=false \
    CI=true

COPY web/ ./
# Use more aggressive optimization for production build
RUN npm run build

FROM golang:1.24-bookworm AS backend-builder

WORKDIR $GOPATH/src/smallest-golang/app/

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd/main.go

FROM gcr.io/distroless/static-debian11

# Set multiple environment variables in a single layer
ENV GIN_MODE=release \
    PLEFI_SERVER__ADDRESS=0.0.0.0:8080 \
    PLEFI_SERVER__STATIC_PATH=/static \
    PLEFI_DATABASE__MIGRATIONS_PATH=/migrations

COPY --from=backend-builder /main .
COPY ./migrations /migrations

# Copy the built frontend assets to the static directory
COPY --from=frontend-builder /app/web/build /static

CMD ["./main"]