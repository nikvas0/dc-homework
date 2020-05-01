# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /src/dc-homework/
COPY notifications/go.mod .
COPY notifications/go.sum .
COPY lib/ ../lib/
RUN go mod download

COPY notifications/ .

RUN GOOS=linux go build -o /go/bin/notifications

FROM golang:latest
COPY --from=builder /go/bin/notifications /go/bin/notifications

CMD ["/go/bin/notifications"]
