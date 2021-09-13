FROM golang:1.17 AS builder
WORKDIR /app
COPY demo-service.go checkavailability.go clientsearch.go config.go parseyandexresponse.go searchsite.go types.go go.mod go.sum ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ds .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/ds .
COPY config.yaml /app/ds
ENTRYPOINT [ "/root/ds" ]
