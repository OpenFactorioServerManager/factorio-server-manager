# Factorio Server Manager Docker Image 
### with SSL reverse-proxy pre-configured

## How to use?

Pull the Docker container from Docker Hub using the pull command
```
docker pull majormjr/factorio-server-manager
```

Now you can start the container by running:

```
docker run --name factorio-manager -d -p 80:80 -p 443:443 -p 34197:34197/udp majormjr/factorio-server-manager
```

If you want persistent data in your container also mount the data volumes when starting:
```
docker run --name factorio-manager -d -v [yourpath]:/opt/factorio/saves [yourpath]:/opt/factorio/mods [yourpath]:/opt/factorio/config -p 80:80 -p 443:443 -p 34197:34197/udp majormjr/factorio-server-manager
```

## Security
This container will auto-generate self-signed ssl certificates on your machine. 

It's advisable that you CHANGE the private key and certificate to one you generated yourself. Do not trust me! Trust yourselves! 

Also if you get an actual SSL certificate for your key it will hide the annoying "the certificate is not trusted blablabla" message.

### But how do I change it?
Nothing easier than that:

When first running the container you need to mount the security volume to your host machine by running `docker run --name factorio-manager -d -v [yourpath]:/security -p 80:80 -p 443:443 -p 34197:34197/udp factorio-server-manager`

The directory will contain a "server.key" file and a "server.crt" file. 

If you replace these with a trusted SSL certificate and key, you should check that "server.crt" contains the whole certificate chain from the root of your CA.

## Updating Credentials, adding and deleting users.
An admin user is created initially using the credentials defined in the factorio-server-manager config file.

The default admin credentials are `user:admin password:factorio`.

Users can be added and deleted on the settings page.

## Updating Factorio
For now you can't update/downgrade the Factorio version from the UI. 

You can however do this using docker images while sustaining your security settings and map/modfiles. 

This guide assumes that you mounted the volumes /security, /opt/factorio/saves, /opt/factorio/config and /opt/factorio/mods to your file system. Before doing anything we need to stop the old container using `docker stop factorio-manager`. To update Factorio you should then open the Dockerfile and change the Factorio version to the one desired. After that you need to rebuild the image using `docker build -t factorio-server-manager .`. Once completed you can simply rerun the command that you used to run the image in the first place. It's recommended to change the name to something including the version to keep track of the containers.

## For everyone who actually read this thing to the end
And now go and build some nice factories!
