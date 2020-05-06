# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /src/dc-homework/
COPY server/go.mod .
COPY server/go.sum .
COPY lib/ ../lib/
RUN go mod download

COPY server/ .

RUN GOOS=linux go build -o /go/bin/shop_server

FROM golang:latest
COPY --from=builder /go/bin/shop_server /go/bin/shop_server

EXPOSE 8080
ENTRYPOINT ["/go/bin/shop_server"]
