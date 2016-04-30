#!/bin/sh
mkdir -p /security
if [ ! -f /security/server.key ]; then
	echo "No SSL key found. generating new key and certificate"
	openssl req \
		-new \
		-newkey rsa:2048 \
		-days 365 \
		-nodes\
		-x509 \
		-subj "/CN=localhost" \
		-keyout /security/server.key \
		-out /security/server.crt
fi

if [ ! -f /security/passwords.conf ]; then
	echo "Generating password file"
	if [ -z "$ADMIN_PASSWORD" ]; then
		echo "Generating credentials"
		export ADMIN_PASSWORD=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c10)
	fi
	echo "Credentials:"
	echo "**********************************"
	echo "Username: admin"
	echo "Password: $ADMIN_PASSWORD"
	echo "**********************************"	
	echo -n "admin:" >> /security/passwords.conf
	openssl passwd -apr1 $ADMIN_PASSWORD >> /security/passwords.conf
fi

nohup nginx &
/opt/factorio-server/factorio-server-manager -dir /opt/factorio
