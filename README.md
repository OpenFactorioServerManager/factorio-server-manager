#Factorio Server Manager

###A tool for managing dedicated Factorio servers.
This tool runs on a Factorio server and allows management of saves, mods and many other features.

![Factorio Server Manager Screenshot](http://i.imgur.com/EbRM03Z.png "Factorio Server Manager")

## Installation
1. Clone the repository
2. Build the binary from the repository root (go build)
3. Run the program

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

```



## Building the server
The backend is built as a REST API via the Go web application.  

It also acts as the webserver to serve the front end react application

All api actions are accessible with the /api route.  The frontend is accessible with the root url.
###Building Go backend
git clone
cd
go build

###building React frontend
install nodejs (use )
cd static/js
npm install
webpack
cp index.html bundle.js dist/ /app/
