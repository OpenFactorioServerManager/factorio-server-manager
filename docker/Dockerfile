# Glibc is required for Factorio Server binaries to run
FROM ubuntu

ENV FACTORIO_VERSION=stable \
    MANAGER_VERSION=0.10.1 \
    RCON_PASS=""

VOLUME /opt/fsm-data /opt/factorio/saves /opt/factorio/mods /opt/factorio/config

EXPOSE 80/tcp 34197/udp

RUN apt-get update && apt-get install -y curl tar xz-utils unzip jq && rm -rf /var/lib/apt/lists/*

WORKDIR /opt

# Install FSM
RUN curl --location "https://github.com/OpenFactorioServerManager/factorio-server-manager/releases/download/$MANAGER_VERSION/factorio-server-manager-linux-$MANAGER_VERSION.zip" \
         --output /tmp/factorio-server-manager-linux_${MANAGER_VERSION}.zip \
    && unzip /tmp/factorio-server-manager-linux_${MANAGER_VERSION}.zip \
    && rm /tmp/factorio-server-manager-linux_${MANAGER_VERSION}.zip \
    && mv factorio-server-manager fsm

COPY ./entrypoint.sh /opt/entrypoint.sh

ENTRYPOINT ["/opt/entrypoint.sh"]
