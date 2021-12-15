# build stage
FROM golang:1.17.5-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.14
WORKDIR /app
COPY /static static
COPY --from=builder /app/main .

EXPOSE 8080
CMD [ "/app/main" ]