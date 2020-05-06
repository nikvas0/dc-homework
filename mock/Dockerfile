# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /src/dc-homework/
COPY mock/go.mod .
COPY mock/go.sum .
COPY lib/ ../lib/
RUN go mod download

COPY mock/ .

RUN GOOS=linux go build -o /go/bin/mock

FROM golang:latest
COPY --from=builder /go/bin/mock /go/bin/mock

CMD ["/go/bin/mock"]
