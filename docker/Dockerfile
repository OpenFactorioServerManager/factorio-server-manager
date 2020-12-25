# Glibc is required for Factorio Server binaries to run
FROM frolvlad/alpine-glibc

ENV FACTORIO_VERSION=stable \
    MANAGER_VERSION=0.8.2 \
    ADMIN_USER=admin \
    ADMIN_PASS=factorio \
    RCON_PASS="" \
    COOKIE_ENCRYPTION_KEY=""

VOLUME /opt/fsm-data /opt/factorio/saves /opt/factorio/mods /opt/factorio/config

EXPOSE 80/tcp 34197/udp

RUN apk add --no-cache curl tar xz unzip jq

WORKDIR /opt

# Install Factorio
RUN curl -s -L -S -k https://www.factorio.com/get-download/$FACTORIO_VERSION/headless/linux64 -o /tmp/factorio_$FACTORIO_VERSION.tar.xz && \
    tar Jxf /tmp/factorio_$FACTORIO_VERSION.tar.xz && \
    rm /tmp/factorio_$FACTORIO_VERSION.tar.xz && \

# Install FSM
RUN curl --location "https://github.com/OpenFactorioServerManager/factorio-server-manager/releases/download/$MANAGER_VERSION/factorio-server-manager-linux-$MANAGER_VERSION.zip" \
         --output /tmp/factorio-server-manager-linux_${MANAGER_VERSION}.zip \
    && unzip /tmp/factorio-server-manager-linux_${MANAGER_VERSION}.zip \
    && rm /tmp/factorio-server-manager-linux_${MANAGER_VERSION}.zip \
    && mv factorio-server-manager fsm

# Setup Nginx
    mkdir -p /run/nginx && \
    chown nginx:root /var/lib/nginx


COPY ./entrypoint.sh /opt/entrypoint.sh

ENTRYPOINT ["/opt/entrypoint.sh"]