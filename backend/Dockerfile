FROM golang:1.26-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev git \
    && go install github.com/air-verse/air@latest

COPY go.mod ./
RUN go mod download

COPY . .

CMD ["air"]