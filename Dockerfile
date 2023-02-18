FROM golang:1.19-alpine as builder

ARG DNS
ARG PORT
ARG SERVER_URL
ARG CASH_HOST
ARG IS_PROD

WORKDIR /shop

COPY . .


ENV DNS=$DNS
ENV PORT=$PORT
ENV SERVER_URL=$SERVER_URL
ARG CASH_HOST=$CASH_HOST
ARG IS_PROD=$IS_PROD

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

EXPOSE 8080

ENTRYPOINT ["/bin/shop"]

