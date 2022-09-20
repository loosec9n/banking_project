#build stage
FROM golang:1.19.1-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#final run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .


EXPOSE 8080

CMD ["/app/main"]