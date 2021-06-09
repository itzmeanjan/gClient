FROM golang:alpine as builder
RUN apk add --no-cache make gcc musl-dev linux-headers git
WORKDIR /app
COPY . .
RUN cd publisher; go build -o ../gclient-pub; cd ../

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/gclient-pub /usr/bin
ENV ADDR="172.17.0.2"
ENV PORT=13000
ENV CLIENT_COUNT=8192
ENV REPEAT_COUNT=0
ENV DELAY="1s"
ENV LOG="false"
ENTRYPOINT [ "gclient-pub", "-topic", "a" ]
