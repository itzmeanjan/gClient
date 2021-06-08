FROM golang:alpine as builder
RUN apk add --no-cache make gcc musl-dev linux-headers git
WORKDIR /app
COPY . .
RUN cd subscriber; go build -o ../gclient-sub; cd ../

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/gclient-sub /usr/bin
ENTRYPOINT [ "gclient-sub", "-addr", "172.17.0.2", "-port", "13000", "-topic", "a", "-out", "false", "-client", "8192" ]
