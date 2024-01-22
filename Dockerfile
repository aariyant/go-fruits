FROM    golang:1.20

WORKDIR /app
RUN    apt update && apt install -y dpkg-dev

COPY    . .

ENV    GOPATH=/usr/local/go/bin/
RUN    pwd && cat .env && go mod tidy
RUN    env GOOS=linux GOARCH=$(dpkg-architecture -q DEB_BUILD_ARCH) go build -v -o app

EXPOSE 8080
CMD    ["/app/app"]

