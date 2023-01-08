FROM golang:1.19-alpine as builder

WORKDIR /shop

COPY . .

RUN go get -d -v ./...

RUN go build -o /bin/shop ./cmd/onlineShopBackend/.

RUN mkdir /bin/static

COPY ./static /bin/static

FROM alpine:latest

RUN mkdir /bin/static

COPY --from=builder /bin/shop /bin

COPY --from=builder /bin/static /bin/static

WORKDIR /bin

VOLUME /bin/static

EXPOSE 8000

ENTRYPOINT ["/bin/shop"]

