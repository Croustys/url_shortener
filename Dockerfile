FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o url_shortener main.go

EXPOSE 3600

CMD [ "./url_shortener" ]