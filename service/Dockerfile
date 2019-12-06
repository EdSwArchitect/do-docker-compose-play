#
# This Dockerfile builds a recent curl with HTTP/2 client support, using
# a recnent nghttp2 build.
#
# See the Makefile for how to tag it. If Docker and that image is found, the
# Go tests use this curl binary for integration tests.
#
FROM ubuntu:18.04 AS EDLOADER3

RUN mkdir -p /opt/playground
RUN mkdir -p /data

RUN ls -l /

WORKDIR /opt/playground

COPY service ed-query

EXPOSE 19090

VOLUME /data

ENTRYPOINT ["/opt/playground/ed-query", "-host", "http://es01:9200"]

#    entrypoint: [ "/opt/playground/ed-query", "-host", "http://es01:9200" ]
