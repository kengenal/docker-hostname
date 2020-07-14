FROM golang:alpine3.11

MAINTAINER kengenal

WORKDIR /app
COPY main.go .

RUN go build -o main .

FROM alpine:3.12

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=0 /app .

EXPOSE 80
CMD ["./main"]
