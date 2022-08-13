FROM alpine:latest

ADD  webhook /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/webhook
