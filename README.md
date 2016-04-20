#Factorio Server Manager

A tool for managing both local and remote Factorio servers.

![Factorio Server Manager Screenshot](http://i.imgur.com/EbRM03Z.png "Factorio Server Manager")


Backend is built as a REST api via the Go application.  It also acts as the webserver to serve the front end react application

All api actions are accessible with the /api route.

To build
Building Go backend
git clone
cd
go build

building React frontend
cd static/js
npm install
webpack
cp index.html bundle.js dist/ /app/
