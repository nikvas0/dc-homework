# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /src/dc-homework/
COPY auth/go.mod .
COPY auth/go.sum .
COPY lib/ ../lib/
RUN go mod download

COPY auth/ .

RUN GOOS=linux go build -o /go/bin/auth

FROM golang:latest
COPY --from=builder /go/bin/auth /go/bin/auth

EXPOSE 8080
ENTRYPOINT ["/go/bin/auth"]
