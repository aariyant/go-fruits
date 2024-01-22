FROM    golang:1.20 as builder

WORKDIR /app
RUN    apt update && apt install -y dpkg-dev

COPY    . .

ENV    GOPATH=/usr/local/go/bin/
RUN    pwd && cat .env && go mod tidy
RUN    env GOOS=linux GOARCH=$(dpkg-architecture -q DEB_BUILD_ARCH) go build -v -o app

FROM    golang:1.20.13-alpine3.19

WORKDIR /app
COPY    --from=builder /app/app /app/.env /app/

EXPOSE 8080
CMD    ["/app/app"]

