FROM golang:1.22-alpine AS build
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . .
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o desafio-tecnico-1

FROM scratch
WORKDIR /app
COPY --from=build /app /usr/bin/server
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app .
EXPOSE 8081
ENTRYPOINT ["./desafio-tecnico-1"]