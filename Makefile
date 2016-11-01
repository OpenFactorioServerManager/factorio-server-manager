# Build tool for Factorio Server Manager
#

build:
	# Build Linux release
	go build -o factorio-server-linux/factorio-server-manager src/*
	ui/node_modules/webpack/bin/webpack.js ui/webpack.config.js app/bundle.js --progress --profile --colors 
	cp -r app/ factorio-server-linux/
	cp conf.json.example factorio-server-linux/conf.json
	zip -r build/factorio-server-linux-x64.zip factorio-server-linux
	rm -rf factorio-server-linux
	# Build Windows release
	GOOS=windows GOARCH=386 go build -o factorio-server-windows/factorio-server-manager.exe src/*
	cp -r app/ factorio-server-windows/
	cp conf.json.example factorio-server-windows/conf.json
	zip -r build/factorio-server-windows.zip factorio-server-windows
	rm -rf factorio-server-windows

