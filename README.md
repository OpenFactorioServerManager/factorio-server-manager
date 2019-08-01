[![Build Status](https://travis-ci.org/mroote/factorio-server-manager.svg?branch=master)](https://travis-ci.org/mroote/factorio-server-manager)

# Factorio Server Manager

### A tool for managing Factorio servers.
This tool runs on a Factorio server and allows management of the Factorio server, saves, mods and many other features.

## Features
* Allows control of the Factorio Server, starting and stopping the Factorio binary.
* Allows the management of save files, upload, download and delete saves.
* Manage installed mods, upload new ones and more
* Manage modpacks, so it is easier to play with different configurations
* Allow viewing of the server logs and current configuration.
* Authentication for protecting against unauthorized users
* Available as a Docker container

## Installation Docker
1. Pull the Docker container from Docker Hub using the pull command
   ```
   docker pull majormjr/factorio-server-manager
   ```

2. Now you can start the container by running:
   ```
   docker run --name factorio-manager -d -p 80:80 -p 443:443 -p 34197:34197/udp majormjr/factorio-server-manager
   ```

## Installation Linux
1. Download the latest release
  * [https://github.com/mroote/factorio-server-manager/releases](https://github.com/mroote/factorio-server-manager/releases)
2. Download the Factorio Standalone server and install to a known directory.
3. Run the server binary file, use the --dir flag to point the management server to your Factorio installation. If you are using the steam installation, point FSM to the steam directory.
  * ```./factorio-server-manager --dir /home/user/.factorio ```
  * ```./factorio-server-manager --dir /home/user/.steam/steam/steamapps/common/Factorio ```
4. Visit [localhost:8080](localhost:8080) in your web browser.

## Installation Windows
1. Download the latest release
  * [https://github.com/mroote/factorio-server-manager/releases](https://github.com/mroote/factorio-server-manager/releases)
2. Download the Factorio Standalone server and install to a known directory.
3. Run the server binary file via cmd or Powershell, use the --dir flag to point the management server to your Factorio installation.
  * ```.\factorio-server-manager --dir C:/Users/username/Factorio```
4. Visit [localhost:8080](localhost:8080) in your web browser.

## Usage
Run the UI server and  specify the directory of your Factorio server installation and the interface to run the HTTP server on.  Edit the conf.json file with your desired credentials for authentication.
```
Usage of ./factorio-server-manager:
  -bin string
    	Location of Factorio Server binary file (default "bin/x64/factorio")
  -conf string
    	Specify location of Factorio Server Manager config file. (default "./conf.json")
  -config string
    	Specify location of Factorio config.ini file (default "config/config.ini")
  -dir string
    	Specify location of Factorio directory. (default "./")
  -host string
    	Specify IP for webserver to listen on. (default "0.0.0.0")
  -max-upload int
    	Maximum filesize for uploaded files (default 20MB). (default 20971520)
  -port string
    	Specify a port for the server. (default "8080")
  -glibc-custom string 
        Specify if custom glibc is used (default false) [true/false]
  -glibc-loc string
        Path to the glibc ld.so file (default "/opt/glibc-2.18/lib/ld-2.18.so")
  -glibc-lib-loc
        Path to the glibc lib folder (default "/opt/glibc-2.18/lib")
        
Example:

./factorio-server-manager --dir /home/user/.factorio --host 10.0.0.1

Custom glibc example:

./factorio-server-manager --dir /home/user/.factorio --host 10.0.0.1 --glibc-custom true --glibc-loc /opt/glibc-2.18/lib/ld-2.18.so --glibc-lib-loc /opt/glibc-2.18/lib

```

## Manage Factorio Server
![Factorio Server Manager Screenshot](http://i.imgur.com/q7tbzdH.png "Factorio Server Manager")

## Manage save files
![Factorio Server Manager Screenshot](http://i.imgur.com/M7kBAhI.png "Factorio Server Manager")

## Manage mods
![Factorio Server Manager Screenshot](https://imgur.com/QIb0Kr4.png "Factorio Server Manager")

## Manage modpacks
![Factorio Server Manager Screenshot](https://imgur.com/O701fB8.png "Factorio Server Manager")



## Development
The backend is built as a REST API via the Go web application.

It also acts as the webserver to serve the front end react application

All api actions are accessible with the /api route.  The frontend is accessible from /.

#### Requirements
+ Go 1.7
+ NodeJS 4.2.6

#### Building Releases
Creates a release zip for windows and linux: (this will install the dependencies listed in gopkgdeps)
```
git clone https://github.com/mroote/factorio-server-manager.git
cd factorio-server-manager
make gen_release
```

#### Building a Testing Binary:
```
git clone https://github.com/mroote/factorio-server-manager.git
cd factorio-server-manager
make
./factorio-server-manager/factorio-server-manager
```

#### Building the React Frontend alone
Frontend is built using React and the AdminLTE CSS framework.

The root of the UI application is served at app/index.html.  Run the npm build script and the Go application during development to get live rebuilding of the UI code.

All necessary CSS and Javascript files are included for running the UI.

Transpiled bundle.js application is output to app/bundle.js, 'npm run build' script starts webpack to build the React application for development.
```
make app/bundle
```

##### For development
The frontend is completly build by npm with laravel-mix. All plugins are buld into the compiled files. No plugins need to be load fro external sources.

It has different variants to build the frontend, provided by laravel-mix:
- `npm run dev` Build the code for development. This will also generate map-files, so the browser, can show, what line and file causes the output.
- `npm run watch` Build the code for development like the dev-command. This will not stop and automatically rebuild, when files are changed and saved.
- `npm run hot` Build the code for development. It has the same behaviour like the watch-command and also causes a hotReload of the files inside the browser (in theory)
- `npm run build` Build the code for deployment. It will generate no map-files and also minifies the bundle-files.
In every of those cases, also images and fonts will be copied to the app-folder.

### Building for Windows
1. Download the latest release source zip file
  * [https://github.com/mroote/factorio-server-manager/releases](https://github.com/mroote/factorio-server-manager/releases)
2. Unzip the Factorio Standalone server and move it to a known directory.
3. Download and install Go 1.7 or newer. https://golang.org/dl/
4. Download and install NodeJS 4.2.6 64-bit or 32-bit depending on your operating system, if unsure download 32-bit
  * https://nodejs.org/download/release/v4.2.6/node-v4.2.6-x64.msi 64-bit
  * https://nodejs.org/download/release/v4.2.6/node-v4.2.6-x86.msi 32-bit
5. Download and install NVM, when asked if you want it to use NodeJS 4.2.6 accept
  * https://github.com/coreybutler/nvm-windows/releases/download/1.1.1/nvm-setup.zip
6. You will need to setup GOPATH in environmental settings in windows. You will want to go into Control Panel\System and Security\System From there on the left hand side click "Advanced system settings". A window will open and you need to click Environment Variables.
7. Under System Variables click New. For Variable name use GOPATH and Variable value C:\Go\

Once everything is installed and ready to go you will need to compile the source for windows

1. Open the folder where ever you unzipped from step #2 above.
2. My folder structure is like this "C:\FS\factorio-server-manager\" C:\FS is where my factorio files are located C:\FS\factorio-server-manager\ is where the server manager files are.
3. You will now need to install some dependencies for Go. You will need to open up a command prompt and one at a time type each of these and hit enter before typing the next one.

```
go get github.com/apexskier/httpauth
go get github.com/go-ini/ini
go get github.com/gorilla/mux
go get github.com/hpcloud/tail
go get github.com/gorilla/websocket
go get github.com/majormjr/rcon
```

3. Now you will want to go into the src folder for example "C:\FS\factorio-server-manager\src" once there hold down left shift and right click an empty area of the folder. Then click "Open command windows here"
4. Type this into the command prompt then hit enter:

```
go build
```

5. Once finished you will now see src.exe or src file inside the folder. You need to move that file to the C:\FS\factorio-server-manager\ or the folder that is before your src folder.
6. From here you need to build the web front-end by going into your ui folder for me its C:\FS\factorio-server-manager\ui\ and again hold shift and left click in an empty area then select open command prompt here. You then need to type this:

```
 npm install
 npm run build
```

7. Now execute the src file created in step #4 above
8. You can now Visit [localhost:8080](localhost:8080) in your web browser to start using the Factorio server Manager

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## Authors

* **Mitch Roote** - [roote.ca](https://roote.ca)

## Special Thanks
- **[All Contributions](https://github.com/mroote/factorio-server-manager/graphs/contributors)**
- **mickael9** for reverseengineering the factorio-save-file: https://forums.factorio.com/viewtopic.php?f=5&t=8568#

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
