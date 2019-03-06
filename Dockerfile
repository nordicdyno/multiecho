FROM golang:1.12.0-alpine3.9 AS builder
WORKDIR /build
ADD . ./
RUN go install .

FROM alpine:3.9

COPY --from=builder /go/bin/multiecho /bin/multiecho
ENTRYPOINT [ "/bin/multiecho" ]
