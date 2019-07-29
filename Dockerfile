FROM golang:1.12.7-alpine3.10 AS builder

RUN apk update && apk add --no-cache git

RUN adduser -D -g '' appuser

WORKDIR $GOPATH/src/github.com/solairerove/linden-honey-go-scraper
COPY . .

RUN go get -d -v

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /go/bin/app /go/bin/app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/go/bin/app"]
