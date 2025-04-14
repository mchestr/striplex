# Dockerfile.distroless
FROM golang:1.24-bookworm AS base

WORKDIR $GOPATH/src/smallest-golang/app/

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM gcr.io/distroless/static-debian11

# Set multiple environment variables in a single layer
ENV GIN_MODE=release \
    PORT=8080

COPY --from=base /main .
# Add the views directory to the container
COPY --from=base /go/src/smallest-golang/app/views /views

CMD ["./main"]