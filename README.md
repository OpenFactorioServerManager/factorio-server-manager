#Factorio Server Manager

###A tool for managing dedicated Factorio servers.
This tool runs on a Factorio server and allows management of saves, mods and many other features.

## Manage save files
![Factorio Server Manager Screenshot](http://i.imgur.com/wZqOuBZ.png "Factorio Server Manager")

## Manage mods
![Factorio Server Manager Screenshot](http://i.imgur.com/45ab48W.png "Factorio Server Manager")

## Installation
1. Clone the repository
  * ```git clone https://github.com/MajorMJR/factorio-server-manager.git```
2. Run the server binary file
  * ```./factorio-server-manager --dir /home/user/.factorio --host 10.0.0.1 ```

## Usage
Run the server and  specify the directory of your Factorio server installation and the interface to run the HTTP server on.
```
Usage of ./factorio-server-manager:
  -bin string
        Location of Factorio Server binary file (default "bin/x64/factorio")
  -config string
        Specify location of Factorio config.ini file (default "config/config.ini")
  -dir string
        Specify location of Factorio directory. (default "./")
  -host string
        Specify IP for webserver to listen on. (default "0.0.0.0")
  -max-upload int
        Maximum filesize for uploaded files. (default 100000)
  -port string
        Specify a port for the server. (default "8080")

Example:

./factorio-server-manager --dir /home/user/.factorio --host 10.0.0.1

```

## Development
The backend is built as a REST API via the Go web application.  

It also acts as the webserver to serve the front end react application

All api actions are accessible with the /api route.  The frontend is accessible from /.

### Requirements
+ Go 1.6
+ NodeJS 4.2.6

#### Building the Go backend
Go Application which manages the Factorio server.

API requests for managing the Factorio server are sent to /api.

The frontend code is served by a HTTP file server running on /.
```
git clone https://github.com/MajorMJR/factorio-server-manager.git
cd factorio-server-manager
go build
```

#### Building the React frontend
Frontend is built using React and the AdminLTE CSS framework. See app/dist/ for AdminLTE included files and license.

The root of the UI application is served at app/index.html.  Run the npm build script and the Go application during development to get live rebuilding of the UI code.

All necessary CSS and Javascript files are included for running the UI.

Transpiled bundle.js application is output to app/bundle.js, 'npm run build' script starts webpack to build the React application for development
```
 install nodejs (use nvm)
 cd ui/
 npm install
 npm run build
 Start factorio-server-manager binary in another terminal
```

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## Authors

* **Mitch Roote** - [roote.me](https://roote.me)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
