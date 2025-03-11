FROM golang:1.23.5-alpine AS dev

WORKDIR /app
COPY . /app/
RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM scratch AS prod

COPY --from=dev go/bin/app /
CMD ["/app"]
