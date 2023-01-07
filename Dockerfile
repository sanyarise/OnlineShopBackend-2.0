FROM golang:1.19-alpine as builder

WORKDIR /shop

COPY . .

RUN go get -d -v ./...

RUN go build -o /bin/shop ./cmd/onlineShopBackend/.

FROM alpine:latest

COPY --from=builder /bin/shop /bin

WORKDIR /bin

RUN mkdir /bin/static
RUN mkdir /bin/static/files
RUN mkdir /bin/static/files/items
RUN mkdir /bin/static/files/categories

VOLUME /bin/static

EXPOSE 8000

ENTRYPOINT ["/bin/shop"]