FROM golang:alpine as builder
RUN apk add --no-cache make gcc musl-dev linux-headers git
WORKDIR /app
COPY . .
RUN cd subscriber; go build -o ../gclient-sub; cd ../

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/gclient-sub /usr/bin
ENV ADDR="172.17.0.2"
ENV PORT=13000
ENV CLIENT_COUNT=8192
ENV LOG="false"
ENV CAPACITY=256
ENTRYPOINT [ "gclient-sub", "-topic", "a" ]
