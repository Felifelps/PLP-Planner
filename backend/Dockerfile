FROM golang:1.26-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev git \
    && go install github.com/air-verse/air@latest \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]