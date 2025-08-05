FROM golang:1.24.5-alpine AS dev

WORKDIR /app
COPY . /app/
RUN go mod download

# install nodejs and npm for maizzle
RUN apk add --update npm

# install maizzle
RUN cd maizzle && \
    npm install && npm run build && cd -

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine:3.21

# Install root certs for TLS validation
RUN apk add --no-cache ca-certificates curl

COPY --from=dev go/bin/app /
CMD ["/app"]

EXPOSE 1025
