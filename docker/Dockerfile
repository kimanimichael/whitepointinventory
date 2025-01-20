FROM golang:1.21-alpine AS build

WORKDIR /src/go/github.com/mike-kimani/whitepointinventory

COPY ../go.mod ./

COPY ../go.sum ./

RUN go mod download

COPY . .

RUN go build -o build/whitepointinventory ./cmd

FROM scratch

WORKDIR /root/

COPY --from=build /src/go/github.com/mike-kimani/whitepointinventory/build/whitepointinventory .

EXPOSE 8080

CMD ["./whitepointinventory"]
