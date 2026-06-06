FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /build/api \
    ./cmd/api/main.go

FROM alpine:3.22.4

RUN apk add --no-cache ca-certificates tzdata

USER nobody

COPY --from=builder /build/api /api

EXPOSE 7001

ENTRYPOINT ["/api"]
