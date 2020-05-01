# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /go/src/github.com/nikvas0/dc-homework/product_upload_inserter
COPY . .

RUN go get -d -v
RUN GOOS=linux go build -o /go/bin/product_upload_inserter

FROM golang:latest
COPY --from=builder /go/bin/product_upload_inserter /go/bin/product_upload_inserter

CMD ["/go/bin/product_upload_inserter"]
