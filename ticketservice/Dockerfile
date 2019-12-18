FROM golang:1.13.0-stretch AS builder
ENV CGO_ENABLED=0
WORKDIR /src
COPY . .
RUN go build -o build/app  app.go

FROM alpine
WORKDIR run
COPY --chown=0:0 --from=builder /src/build /run
ENTRYPOINT ["./app"]