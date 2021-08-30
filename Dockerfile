FROM alpine:latest

RUN mkdir /app
WORKDIR /app
ADD demo-service  /app/demo-service

CMD ["./demo-service"]
