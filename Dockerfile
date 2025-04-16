# Multi-stage build for both frontend and backend
FROM node:22-slim AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

FROM golang:1.24-bookworm AS backend-builder

WORKDIR $GOPATH/src/smallest-golang/app/

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM gcr.io/distroless/static-debian11

# Set multiple environment variables in a single layer
ENV GIN_MODE=release \
    PLEFI_SERVER__ADDRESS=0.0.0.0:8080 \
    PLEFI_SERVER__STATIC_PATH=/static

COPY --from=backend-builder /main .

# Copy the built frontend assets to the static directory
COPY --from=frontend-builder /app/frontend/build /static

CMD ["./main"]