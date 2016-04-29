# Docker Image with SSL reverse-proxy and authentication 

## How to use?
First run `docker build -t factorio-server-manager` in this directory.

Now you can start the container by running `docker run --name factorio-manager -d -p80:80 -p443:443 -p34197:34197/udp factorio-server-manager`

Your default credentials can be retrieved by checking the output of `docker logs factorio-manager`

Ok, with that out of the way, lets talk about security:

## Security. This is important!

I have done my best to secure the container pretty well. This includes the generation of self-signed ssl certificates on your machine. However, I HIGHLY ADVISE you that you CHANGE the private key and certificate to one you generated yourself. Do not trust me! Trust yourselves! Also if you get an actual SSL certificate for your key it will hide the annoying "the certificate is not trusted blablabla" message.

### But how do I change it?
Nothing easier than that:

When first running the container you need to mount the security volume to your host machine by running `docker run --name factorio-manager -d -v [yourpath]:/security -p80:80 -p443:443 -p34197:34197/udp factorio-server-manager`

You should always do that, as this will allow you to change the login credentials for any users as well. The directory will contain a "server.key" file and a "server.crt" file. If you replace these with a trusted SSL certificate and key, you should check that "server.crt" contains the whole certificate chain from the root of your CA.

Ok, you got me. There might be things that are easier than that... You should do it anyways. 

## Updating Credentials, adding and deleting users.
This is where I got lazy. I'm sorry, but I did not create a great tool that automagically does everything for you. But you can do it. As I'm sure you've read the security chapter and you've done everything I said there you should've mounted the security volume to any point on your filetree already. If not, read the security chapter!

In the mounted security volume you'll find a passwords.conf file. This contains encrypted passwords for every user who can access the manager. The format is `username:encryptedpassword`.

Deleting users is pretty straightforward. Delete the correct line. 

To create a new password entry, you can use `openssl passwd -apr1 yourpasswordhere`. That should get you started. 

## For everyone who actually read this thing to the end
You can also set your default admin password by passing it to your initial docker run command like this:

`docker run -d --name factorio-manager -d -v [yourpath]:/security -p80:80 -p443:443 -p34197:34197/udp -e "ADMIN_PASSWORD=jqkSnQS4rA" factorio-server-manager`

And now go and build some nice factories# Docker Image with SSL reverse-proxy and authentication 

## How to use?
First run `docker build -t factorio-server-manager` in this directory.

Now you can start the container by running `docker run --name factorio-manager -d -p80:80 -p443:443 -p34197:34197/udp factorio-server-manager`

Your default credentials can be retrieved by checking the output of `docker logs factorio-manager`

Ok, with that out of the way, lets talk about security:

## Security. This is important!

I have done my best to secure the container pretty well. This includes the generation of self-signed ssl certificates on your machine. However, I HIGHLY ADVISE you that you CHANGE the private key and certificate to one you generated yourself. Do not trust me! Trust yourselves! Also if you get an actual SSL certificate for your key it will hide the annoying "the certificate is not trusted blablabla" message.

### But how do I change it?
Nothing easier than that:

When first running the container you need to mount the security volume to your host machine by running `docker run --name factorio-manager -d -v [yourpath]:/security -p80:80 -p443:443 -p34197:34197/udp factorio-server-manager`

You should always do that, as this will allow you to change the login credentials for any users as well. The directory will contain a "server.key" file and a "server.crt" file. If you replace these with a trusted SSL certificate and key, you should check that "server.crt" contains the whole certificate chain from the root of your CA.

Ok, you got me. There might be things that are easier than that... You should do it anyways. 

## Updating Credentials, adding and deleting users.
This is where I got lazy. I'm sorry, but I did not create a great tool that automagically does everything for you. But you can do it. As I'm sure you've read the security chapter and you've done everything I said there you should've mounted the security volume to any point on your filetree already. If not, read the security chapter!

In the mounted security volume you'll find a passwords.conf file. This contains encrypted passwords for every user who can access the manager. The format is `username:encryptedpassword`.

Deleting users is pretty straightforward. Delete the correct line. 

To create a new password entry, you can use `openssl passwd -apr1 yourpasswordhere`. That should get you started. 

## For everyone who actually read this thing to the end
You can also set your default admin password by passing it to your initial docker run command like this:

`docker run -d --name factorio-manager -d -v [yourpath]:/security -p80:80 -p443:443 -p34197:34197/udp -e "ADMIN_PASSWORD=jqkSnQS4rA" factorio-server-manager`

And now go and build some nice factories!
