FROM golang:1.13.6-alpine AS builder

WORKDIR /go/src/RbacSampleApp

COPY . .
RUN go build -o /RbacSampleApp .

FROM alpine:3.11.2 AS release
RUN apk add --no-cache ca-certificates

WORKDIR /RbacSampleApp
COPY --from=builder /RbacSampleApp ./server
EXPOSE 3550
ENTRYPOINT ["/RbacSampleApp/server"]
