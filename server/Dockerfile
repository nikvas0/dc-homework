# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /go/src/github.com/nikvas0/dc-homework/server
COPY . .

RUN go get -d -v
RUN GOOS=linux go build -o /go/bin/shop_server

FROM golang:latest
COPY --from=builder /go/bin/shop_server /go/bin/shop_server
COPY ./tools/wait.sh /tools/wait.sh
RUN chmod 777 /tools/wait.sh

EXPOSE 8080
ENTRYPOINT ["/go/bin/shop_server"]
