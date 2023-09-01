FROM golang:1.21.0 as builder

COPY . /build

WORKDIR /build

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o main .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/main .

EXPOSE 5011

CMD ["./main"] 
