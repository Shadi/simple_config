FROM golang:alpine as builder

WORKDIR /app 

COPY . .

RUN apk --no-cache add ca-certificates

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/simple_config /usr/bin/

EXPOSE 8080

ENV PORT=8080
ENV GIN_MODE="release"

ENTRYPOINT ["simple_config"]
