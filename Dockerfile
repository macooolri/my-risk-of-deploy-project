FROM golang:1.20rc1-alpine3.17 as builder
WORKDIR /app
COPY /go.mod .
COPY /go.sum .
RUN go mod tidy
COPY . .
RUN go build -o main .

FROM alpine:latest as runner
COPY --from=builder /app/main .
CMD ["./main"]

