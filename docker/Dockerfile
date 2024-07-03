FROM golang:1.21-alpine AS build

WORKDIR /src/go/github.com/mike-kimani/whitepointinventory

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o whitepointinventory

FROM alpine:latest

WORKDIR /root/

COPY --from=build /src/go/github.com/mike-kimani/whitepointinventory/whitepointinventory .

EXPOSE 8080

CMD ["./whitepointinventory"]


