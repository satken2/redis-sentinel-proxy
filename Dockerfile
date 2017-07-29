FROM alpine:latest

MAINTAINER Jérémie BORDIER <jeremie.bordier@gmail.com>

# copy binary
COPY redis-sentinel-proxy /usr/local/bin/redis-sentinel-proxy
COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["redis-sentinel-proxy"]