#Factorio Server Manager

###A tool for managing dedicated Factorio servers.
This tool runs on a Factorio server and allows management of saves, mods and many other features.

![Factorio Server Manager Screenshot](http://i.imgur.com/EbRM03Z.png "Factorio Server Manager")

## Installation
1. Clone the repository
  * ```git clone https://github.com/MajorMJR/factorio-server-manager.git```
2. Run the server binary file
  * ```./factorio-server-manager --dir /home/user/.factorio --host 10.0.0.1 ```

## Usage
Run the server and  specify the directory of your Factorio server installation and the interface to run the HTTP server on.
```
Usage of ./factorio-server-manager:
  -dir string
        Specify location of Factorio config directory. (default "./")
  -host string
        Specify IP for webserver to listen on. (default "0.0.0.0")
  -port string
        Specify a port for the server (default "8080")

Example:

./factorio-server-manager --dir /home/user/.factorio --host 10.0.0.1

```

## Building the server
The backend is built as a REST API via the Go web application.  

It also acts as the webserver to serve the front end react application

All api actions are accessible with the /api route.  The frontend is accessible with the root url.

### Requirements
+ Go 1.6
+ NodeJS 4.2.6

###Building Go backend
git clone https://github.com/MajorMJR/factorio-server-manager.git
cd factorio-server-manager
go build

###building React frontend
install nodejs (use nvm)
cd static/js
npm install
npm run build

## Authors

* **Mitch Roote** - [roote.me](https://roote.me)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
