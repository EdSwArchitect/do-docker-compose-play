#
# This Dockerfile builds a recent curl with HTTP/2 client support, using
# a recnent nghttp2 build.
#
# See the Makefile for how to tag it. If Docker and that image is found, the
# Go tests use this curl binary for integration tests.
#
FROM ubuntu:18.04 AS EDLOADER2

RUN mkdir -p /opt/playground
RUN mkdir -p /data

RUN ls -l /

WORKDIR /opt/playground

COPY loader ed-loader

EXPOSE 18080

VOLUME /data

ENTRYPOINT ["/opt/playground/ed-loader"]
