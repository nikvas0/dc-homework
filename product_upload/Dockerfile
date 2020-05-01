# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /src/dc-homework/
COPY product_upload/go.mod .
COPY product_upload/go.sum .
COPY lib/ ../lib/
RUN go mod download

COPY product_upload/ .

RUN GOOS=linux go build -o /go/bin/product_upload

FROM golang:latest
COPY --from=builder /go/bin/product_upload /go/bin/product_upload

EXPOSE 8080
ENTRYPOINT ["/go/bin/product_upload"]
