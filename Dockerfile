FROM golang:1.19 as builder

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -o server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/app/server .
COPY --from=builder /go/src/app/assets assets
COPY --from=builder /go/src/app/public public

EXPOSE 80

CMD ["./server", "--hostport", ":80"]


