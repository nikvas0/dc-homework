# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /go/src/github.com/nikvas0/dc-homework
COPY . .
WORKDIR /go/src/github.com/nikvas0/dc-homework/server

RUN go get -d -v
RUN GOOS=linux go build -o /go/bin/shop_server

FROM golang:latest
COPY --from=builder /go/bin/shop_server /go/bin/shop_server

EXPOSE 8080
ENTRYPOINT ["/go/bin/shop_server"]