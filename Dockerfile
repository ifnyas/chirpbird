FROM golang:1.17.5-alpine3.14
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 8080
CMD [ "/app/main" ]