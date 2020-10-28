# Factorio Server Manager Docker Image

## Prerequisites
You need to have [Docker](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-20-04)
and [Docker Compose](https://docs.docker.com/compose/install/) installed.

## Getting started?

Copy `docker-compose.yaml` and `.env` files from this repository to somewhere on your server.

Edit values in the `.env` file:
* `ADMIN_USER` (default `admin`): Name of the default user created for FSM UI.
* `ADMIN_PASS` (default `factorio`): Default user password. \
  __Important:__ _For security reasons, please change the default user name and password. Never use the defaults._
* `RCON_PASS` (default empty string): Password for Factorio RCON (FSM uses it to communicate with the Factorio server). \
  If left empty, a random password will be generated and saved on the first start of the server. You can see the password in `fsm-data/conf.json` file.
* `COOKIE_ENCRYPTION_KEY` (default empty string): The key used to encrypt auth cookie for FSM UI. \
  If left empty, a random key will be generated and saved on the first start of the server. You can see the key in `fsm-data/conf.json` file.
* `DOMAIN_NAME` (must be set manually): The domain name where your FSM UI will be available. Must be set,
  so [Let's Encrypt](https://letsencrypt.org/) service can issue a valid HTTPS certificate for this domain.
* `EMAIL_ADDRESS` (must be set manually): Your email address. Used only by Let's Encrypt service.

Alternatively you can ignore `.env` file and edit this values directly in `environment` section of `docker-compose.yaml`.
But remember that if `.env` file is present, values set there take precedence over values set in `docker-compose.yaml`.

Now you can start the container by running:

```
docker-compose up -d
```

### Simple configuration without HTTPS

If you don't care about HTTPS and want to run just the Factorio Server Manager, or want to run it on local machine you can use `docker-compose.simple.yaml`.

Ignore `DOMAIN_NAME` and `EMAIL_ADDREESS` variables in `.env` file and run
```
docker-compose -f docker-compose.simple.yaml up -d
```

### Factorio version

By default container will download the latest version of factorio. If you want to use specific version, you can change
the value of `FACTORIO_VERSION=latest` variable in the `docker-compose.yaml` file.

## Accessing the application

Go to the domain specified in your `.env` file in your web browser. If running on localhost host access the application at http://localhost

### First start

When container starts it begins to dowload Factorio headless server archive, and only after that Factorio Server Manager server starts.
So when Docker Compose writes
```
Creating factorio-server-manager ... done
```
you have to wait several seconds before FSM UI becomes available.

It may take some time for Let's Encrypt to issue the certificate, so for the first couple of minutes after starting the container you may see
"Your connection is not private" error when you open your Factorio Server Manager address in your browser. This error should disappear within
a couple of minutes, if configuration parameters are set correctly.

## Updating Credentials, adding and deleting users.

An admin user is created initially using the credentials defined in the factorio-server-manager config file.

Users can be added and deleted on the settings page.

## Updating Factorio

For now you can't update/downgrade the Factorio version from the UI.

You can however do this using docker images while sustaining your security settings and map/modfiles.

If you want to update Factorio to the latest version:
1. Save your game and stop Factorio server in FSM UI.
2. Run `docker-compose restart` (or `docker-compose -f docker-compose.simple.yaml restart` if you are using simple configuration).

After container starts, latest Factorio version will be downloaded and installed.

## Security

Authentication is supported in the application but it is recommended to ensure access to the Factorio manager UI is accessible via VPN or internal network.

## Development
For development purposes it also has the ability to create the docker image from local sourcecode. This is done by running `build.sh` in the `docker` directory. This will delete all old executables and the node_modules directory (runs `make build`). The created docker image will have the tag `factorio-server-manager:dev`.

## For everyone who actually read this thing to the end

And now go and build some nice factories!
