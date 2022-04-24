FROM golang AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine

COPY --from=builder /build/main /

#EXPOSE 8080

#ENTRYPOINT ["/main"]