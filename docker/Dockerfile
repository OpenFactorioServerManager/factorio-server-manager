# glibc is required for Factorio Server binaries to run
FROM frolvlad/alpine-glibc

MAINTAINER Jannik Kolodziej <docker@jkolodziej.de>

ENV FACTORIO_VERSION=0.12.32 \
    MANAGER_VERSION=0.2.0 \
    ADMIN_PASSWORD=

VOLUME /opt/factorio/saves /opt/factorio/mods /security

RUN apk add --no-cache curl tar unzip nginx

WORKDIR /opt/
ADD "init.sh" "/opt/init.sh"
ADD "nginx.conf" "/etc/nginx/nginx.conf"

RUN curl -s -L -S -k https://www.factorio.com/get-download/$FACTORIO_VERSION/headless/linux64 -o /tmp/factorio_$FACTORIO_VERSION.tar.gz && \
    tar zxf /tmp/factorio_$FACTORIO_VERSION.tar.gz && \
    rm /tmp/factorio_$FACTORIO_VERSION.tar.gz && \
    curl -s -L -S -k https://github.com/MajorMJR/factorio-server-manager/releases/download/$MANAGER_VERSION/factorio-server-manager-linux-x64.zip --cacert /opt/github.pem -o /tmp/factorio-server-manager-linux-x64_$MANAGER_VERSION.zip && \
    unzip -qq /tmp/factorio-server-manager-linux-x64_$MANAGER_VERSION.zip && \
    rm /tmp/factorio-server-manager-linux-x64_$MANAGER_VERSION.zip && \
    rm -r /var/lib/nginx && \
    mkdir -p /var/lib/nginx/tmp && \
    chown nginx:root /var/lib/nginx

EXPOSE 80/tcp 443/tcp 34190-34200/udp

ENTRYPOINT ["/opt/init.sh"]
