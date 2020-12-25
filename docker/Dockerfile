# glibc is required for Factorio Server binaries to run
FROM frolvlad/alpine-glibc

ENV FACTORIO_VERSION=stable \
    MANAGER_VERSION=0.8.2 \
    ADMIN_PASSWORD=factorio

VOLUME /opt/factorio/saves /opt/factorio/mods /opt/factorio/config /security

RUN apk add --no-cache curl tar unzip nginx openssl xz

WORKDIR /opt/

RUN curl -s -L -S -k https://www.factorio.com/get-download/$FACTORIO_VERSION/headless/linux64 -o /tmp/factorio_$FACTORIO_VERSION.tar.xz && \
    tar Jxf /tmp/factorio_$FACTORIO_VERSION.tar.xz && \
    rm /tmp/factorio_$FACTORIO_VERSION.tar.xz && \
    curl -sLSk https://github.com/mroote/factorio-server-manager/releases/download/$MANAGER_VERSION/factorio-server-manager-linux-${MANAGER_VERSION}.zip \
     --cacert /opt/github.pem -o /tmp/factorio-server-manager-linux_$MANAGER_VERSION.zip && \
    unzip -qq /tmp/factorio-server-manager-linux_$MANAGER_VERSION.zip && \
    rm /tmp/factorio-server-manager-linux_$MANAGER_VERSION.zip && \
    mkdir -p /run/nginx && \
    chown nginx:root /var/lib/nginx

COPY "init.sh" "/opt/init.sh"
COPY "nginx.conf" "/etc/nginx/nginx.conf"

EXPOSE 80/tcp 443/tcp 34190-34200/udp

ENTRYPOINT ["/opt/init.sh"]
