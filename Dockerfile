FROM golang:1.12
WORKDIR /src/unload
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" .

FROM alpine:3.9
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
RUN apk --no-cache add ca-certificates
COPY --from=0 /src/unload/unload /usr/local/bin
EXPOSE 50051
CMD ["unload"]