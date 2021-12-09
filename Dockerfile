FROM golang:1-alpine as builder

ARG VERSION

RUN go install github.com/korylprince/fileenv@v1.1.0
RUN go install "github.com/korylprince/hasura-ad-webhook@$VERSION"


FROM alpine:3.15

RUN apk add --no-cache bash ca-certificates

COPY --from=builder /go/bin/fileenv /
COPY --from=builder /go/bin/hasura-ad-webhook /
COPY ./setenv.sh /

CMD ["/fileenv", "/bin/bash", "/setenv.sh", "/hasura-ad-webhook"]
