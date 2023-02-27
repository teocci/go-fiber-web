FROM golang:1.17.11 AS builder
RUN apt -y update && apt -y install git
WORKDIR /go/src
COPY . ./

RUN go mod download
RUN go build -o webserver

FROM ubuntu

WORKDIR /app

COPY --from=builder /go/src/webserver ./
COPY --from=builder /go/src/web ./web

COPY --from=builder /go/src/config.json ./

EXPOSE 9012

CMD ["./webserver"]