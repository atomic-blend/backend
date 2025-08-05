FROM golang:1.24.5-alpine AS dev

WORKDIR /app
COPY . /app/
RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine:3.21

# Install root certs for TLS validation
RUN apk add --no-cache ca-certificates curl

COPY --from=dev go/bin/app /
CMD ["/app"]

EXPOSE 1025
