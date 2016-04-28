FROM debian:jessie

ENV FSM_VERSION 0.2.0

RUN apt-get update && \
    apt-get -y install unzip wget ca-certificates && \
    apt-get clean

RUN wget https://github.com/MajorMJR/factorio-server-manager/releases/download/$FSM_VERSION/factorio-server-manager-linux-x64.zip -P /tmp/
RUN unzip /tmp/factorio-server-manager-linux-x64.zip -d /opt/

WORKDIR /opt/factorio-server

VOLUME ["/factorio"]

#Expose UI Porthttps://github.com/MajorMJR/factorio-server-manager/releases/download/$FSM_VERSION/factorio-server-manager-linux-x64.zip
EXPOSE 8080

ENTRYPOINT ["/opt/factorio-server/factorio-server-manager", "--dir", "/factorio"]
