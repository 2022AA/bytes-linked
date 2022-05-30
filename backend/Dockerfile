FROM golang:1.17.6 as builder

COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go install

FROM debian:11.2 as bytes-linked
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /go/bin/bytes-linked /app
COPY conf /app/conf
ENTRYPOINT ["/app/bytes-linked"]
