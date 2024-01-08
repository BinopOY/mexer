FROM golang:latest

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /main
EXPOSE 3000
ENV PORT=3000

# Run
CMD ["/main"]
