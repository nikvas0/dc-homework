# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /src/dc-homework/
COPY product_upload_inserter/go.mod .
COPY product_upload_inserter/go.sum .
COPY lib/ ../lib/
RUN go mod download

COPY product_upload_inserter/ .

RUN GOOS=linux go build -o /go/bin/product_upload_inserter

FROM golang:latest
COPY --from=builder /go/bin/product_upload_inserter /go/bin/product_upload_inserter

CMD ["/go/bin/product_upload_inserter"]
