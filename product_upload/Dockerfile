# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /go/src/github.com/nikvas0/dc-homework/product_upload
COPY . .

RUN go get -d -v
RUN GOOS=linux go build -o /go/bin/auth

FROM golang:latest
COPY --from=builder /go/bin/auth /go/bin/auth

EXPOSE 8080
ENTRYPOINT ["/go/bin/auth"]
