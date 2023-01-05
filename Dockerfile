FROM golang:1.19-alpine as builder

WORKDIR /shop

COPY . .

RUN go get -d -v ./...

RUN go build -o /bin/shop ./cmd/onlineShopBackend/.

FROM alpine:latest

COPY --from=builder /bin/shop /bin

#COPY website /app/website/

WORKDIR /bin

ENTRYPOINT ["/bin/shop"]
