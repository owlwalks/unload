FROM golang:1.13
WORKDIR /src/unload
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" .

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
COPY --from=0 /src/unload/unload /usr/local/bin
EXPOSE 50051
ENTRYPOINT ["unload"]