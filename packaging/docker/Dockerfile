FROM alpine:latest as runtime
MAINTAINER Marc Sauter <marc.sauter@bluewin.ch>

COPY gitlabctl /usr/bin/gitlabctl

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534