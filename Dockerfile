FROM ubuntu:latest

MAINTAINER https://github.com/abliqo

RUN apt-get update && apt-get install -y runit

RUN mkdir -p /opt/mtls-echo/config && \
    mkdir -p /opt/mtls-echo/content && \
    mkdir -p /etc/service/mtls-echo

COPY bin/mtls-echo /opt/mtls-echo/
COPY content/*.* /opt/mtls-echo/content/
COPY config/*.* /opt/mtls-echo/config/
COPY docker/run /etc/service/mtls-echo/
COPY docker/runsvinit.sh /bin/

WORKDIR /opt/mtls-echo

ENTRYPOINT ["/bin/runsvinit.sh"]