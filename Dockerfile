FROM golang:1.25-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o justsruput main.go

FROM alpine:latest 

WORKDIR /backend

COPY --from=builder /app/justsruput /backend/justsruput

EXPOSE 8080

ENTRYPOINT [ "./justsruput" ]