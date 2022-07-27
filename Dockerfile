FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY src/*.go ./

RUN go build -ldflags '-w -extldflags "-static"' -tags netgo -o ipsync

FROM scratch

WORKDIR /app
COPY --from=builder /app/ipsync ./ipsync

CMD ["/app/ipsync"]

