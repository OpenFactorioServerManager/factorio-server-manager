#!/bin/sh

init_config() {
    jq_cmd='.'

    if [ -n $ADMIN_USER ]; then
        jq_cmd="${jq_cmd} | .username = \"$ADMIN_USER\""
        echo "Admin username is '$ADMIN_USER'"
    fi
    if [ -n $ADMIN_PASS ]; then
        jq_cmd="${jq_cmd} | .password = \"$ADMIN_PASS\""
        echo "Admin password is '$ADMIN_PASS'"
    fi
    echo "IMPORTANT! Please create new user and delete default admin user ASAP."

    if [ -z $RCON_PASS ]; then
        RCON_PASS="$(random_pass)"
    fi
    jq_cmd="${jq_cmd} | .rcon_pass = \"$RCON_PASS\""
    echo "Factorio rcon password is '$RCON_PASS'"

    if [ -z $COOKIE_ENCRYPTION_KEY ]; then
        COOKIE_ENCRYPTION_KEY="$(random_pass)"
    fi
    jq_cmd="${jq_cmd} | .cookie_encryption_key = \"$COOKIE_ENCRYPTION_KEY\""

    jq_cmd="${jq_cmd} | .database_file = \"/opt/fsm-data/auth.leveldb\""
    jq_cmd="${jq_cmd} | .log_file = \"/opt/fsm-data/factorio-server-manager.log\""

    jq "${jq_cmd}" /opt/fsm/conf.json >/opt/fsm-data/conf.json
}

random_pass() {
    LC_ALL=C tr -dc 'a-zA-Z0-9' </dev/urandom | fold -w 24 | head -n 1
}

install_game() {
    curl --location "https://www.factorio.com/get-download/${FACTORIO_VERSION}/headless/linux64" \
         --output /tmp/factorio_${FACTORIO_VERSION}.tar.xz \
    && tar -xf /tmp/factorio_${FACTORIO_VERSION}.tar.xz \
    && rm /tmp/factorio_${FACTORIO_VERSION}.tar.xz
}

if [ ! -f /opt/fsm-data/conf.json ]; then
    init_config
fi

install_game

cd /opt/fsm && ./factorio-server-manager -conf /opt/fsm-data/conf.json -dir /opt/factorio -port 80

