# mexer does not work in alpine environment
FROM golang:1.21.5 AS base
FROM base AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main

FROM base AS final

WORKDIR /app

# Copy binary
COPY --from=builder /app/main .
COPY --from=builder /app/views ./views
COPY --from=builder /app/static ./static

# Run
EXPOSE 3000
ENV PORT=3000
CMD ["./main"]
