# Build tool for Factorio Server Manager
#

NODE_ENV:=production

build:
	# Build Linux release
	mkdir build
	GOOS=linux GOARCH=amd64 go build -o factorio-server-linux/factorio-server-manager src/*
#	ui/node_modules/webpack/bin/webpack.js ui/webpack.config.js app/bundle.js --progress --profile --colors 
	cp -r app/ factorio-server-linux/
	cp conf.json.example factorio-server-linux/conf.json
	zip -r build/factorio-server-manager-linux-x64.zip factorio-server-linux
	rm -rf factorio-server-linux
	# Build Windows release
	GOOS=windows GOARCH=386 go build -o factorio-server-windows/factorio-server-manager.exe src/*
	cp -r app/ factorio-server-windows/
	cp conf.json.example factorio-server-windows/conf.json
	zip -r build/factorio-server-manager-windows.zip factorio-server-windows
	rm -rf factorio-server-windows

dev:
	mkdir dev
	GOOS=linux GOARCH=amd64 go build -o factorio-server-linux/factorio-server-manager src/*
	cp -r app/ dev/
	cp conf.json.example dev/conf.json
	mv factorio-server-linux/factorio-server-manager dev/factorio-server-manager
