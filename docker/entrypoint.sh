#!/bin/sh

init_config() {
    jq_cmd='.'

    if [ -n "$RCON_PASS" ]; then
      jq_cmd="${jq_cmd} | .rcon_pass = \"$RCON_PASS\""
      echo "Factorio rcon password is '$RCON_PASS'"
    fi

    jq_cmd="${jq_cmd} | .sq_lite_database_file = \"/opt/fsm-data/sqlite.db\""
    jq_cmd="${jq_cmd} | .log_file = \"/opt/fsm-data/factorio-server-manager.log\""

    jq "${jq_cmd}" /opt/fsm/conf.json >/opt/fsm-data/conf.json
}

random_pass() {
    LC_ALL=C tr -dc 'a-zA-Z0-9' </dev/urandom | fold -w 24 | head -n 1
}

install_game() {
    curl --location "https://www.factorio.com/get-download/${FACTORIO_VERSION}/headless/linux64" \
         --output /tmp/factorio_${FACTORIO_VERSION}.tar.xz
    tar -xf /tmp/factorio_${FACTORIO_VERSION}.tar.xz
    rm /tmp/factorio_${FACTORIO_VERSION}.tar.xz
}

if [ ! -f /opt/fsm-data/conf.json ]; then
    init_config
fi

install_game

cd /opt/fsm && ./factorio-server-manager --conf /opt/fsm-data/conf.json --dir /opt/factorio --port 80 "$@"

