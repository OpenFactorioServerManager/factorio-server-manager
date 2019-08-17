# Factorio Server Manager Docker Image

## Getting started?

Pull the Docker container from Docker Hub using the pull command

```
docker pull majormjr/factorio-server-manager
```

Now you can start the container by running:

```
docker run --name factorio-manager -d \
    -p 80:80 \
    -p 443:443 \
    -p 34197:34197/udp \
    majormjr/factorio-server-manager
```

If you want persistent data in your container also mount the data volumes when starting:

```
docker run --name factorio-manager -d \
    -v [yourpath]:/opt/factorio/saves \
    -v [yourpath]:/opt/factorio/mods \
    -v [yourpath]:/opt/factorio/config \
    -v [yourpath]:/security \
    -p 80:80 \
    -p 443:443 \
    -p 34197:34197/udp \
    majormjr/factorio-server-manager
```


You can also user a docker-compose file.

* Create a file `docker-compose.yml`
* Enter the following data

```
version: 3
services:
    factorio-server-manager:
        container_name: factorio-manager
        volumes:
            - '[yourpath]/saves:/opt/factorio/saves'
            - '[yourpath]/mods:/opt/factorio/mods'
            - '[yourpath]/config:/opt/factorio/config'
            - '[yourpath]/security:/security'
        ports:
            - '80:80'
            - '443:443'
            - '34197:34197/udp'
        image: majormjr/factorio-server-manager
```

Run the file with `docker-compose up`


## Accessing the application

Go to the port specified in your `docker run` command in your web browser. If running on localhost host access the application at https://localhost

## Updating Credentials, adding and deleting users.

An admin user is created initially using the credentials defined in the factorio-server-manager config file.

Users can be added and deleted on the settings page.

## Updating Factorio

For now you can't update/downgrade the Factorio version from the UI.

You can however do this using docker images while sustaining your security settings and map/modfiles.

This guide assumes that you mounted the volumes /security, /opt/factorio/saves, /opt/factorio/config and /opt/factorio/mods to your file system. Before doing anything we need to stop the old container using `docker stop factorio-manager`. To update Factorio you should then open the Dockerfile and change the Factorio version to the one desired. After that you need to rebuild the image using `docker build -t factorio-server-manager .`. Once completed you can simply rerun the command that you used to run the image in the first place. It's recommended to change the name to something including the version to keep track of the containers.

Pull the latest container with `docker pull majormjr/factorio-server-manager` and start with the `docker run` command.

## Security

A self generated SSL/TLS certificate is created when the container is first created and the webserver is accessible via HTTPS.

Authentication is supported in the application but it is recommended to ensure access to the Factorio manager UI is accessible via VPN or internal network.

### Changing SSL/TLS certificate

If you have your own SSL/TLS certificate then you can supply it to the Factorio Server Manager container.

When first running the container you need to mount the security volume to your host machine by adding the security volume parameter `-v [yourpath]:/security`

The directory will contain a "server.key" file and a "server.crt" file.

If you replace these with a trusted SSL certificate and key, you should ensure that "server.crt" contains the whole certificate chain from the root of your CA.

## For everyone who actually read this thing to the end

And now go and build some nice factories!
