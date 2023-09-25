FROM golang:1.19.13-alpine as builder

COPY src /ipsync
WORKDIR /ipsync
RUN CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"'

FROM scratch

COPY --from=builder /ipsync/ipsync /sbin/ipsync
CMD [ "/sbin/ipsync" ]
