# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /go/src/github.com/nikvas0/dc-homework/notifications
COPY . .

RUN go get -d -v
RUN GOOS=linux go build -o /go/bin/notifications

FROM golang:latest
COPY --from=builder /go/bin/notifications /go/bin/notifications

CMD ["/go/bin/notifications"]
