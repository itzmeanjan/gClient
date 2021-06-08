FROM golang:alpine as builder
RUN apk add --no-cache make gcc musl-dev linux-headers git
WORKDIR /app
COPY . .
RUN cd publisher; go build -o ../gclient-pub; cd ../

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/gclient-pub /usr/bin
ENTRYPOINT [ "gclient-pub", "-addr", "172.17.0.2", "-port", "13000", "-topic", "a", "-topic", "b", "-repeat", "0", "-client", "1", "-delay", "1s", "-out", "false" ]
