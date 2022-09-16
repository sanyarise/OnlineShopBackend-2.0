FROM golang:1.19-alpine as builder

WORKDIR .

COPY . .

RUN go get -d -v ./...

RUN go build -o /bin/shop ./cmd/onlineShopBackend/.

FROM alpine:latest

COPY --from=builder /bin/url /bin

#COPY website /app/website/

WORKDIR /bin

ENTRYPOINT ["/bin/shop"]