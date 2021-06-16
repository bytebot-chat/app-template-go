FROM golang:1.16.4-alpine3.13 as builder

RUN adduser -D -g 'bytebot' bytebot
WORKDIR /app
COPY . .
RUN ./docker-build.sh

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/opt/app /opt/app
VOLUME /data

USER bytebot
ENTRYPOINT ["/opt/app"]
