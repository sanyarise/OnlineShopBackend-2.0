FROM golang:1.10 AS build
WORKDIR /go/src
COPY api ./api
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o onlineShopBackend .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/onlineShopBackend ./
EXPOSE 8080/tcp
ENTRYPOINT ["./onlineShopBackend"]
