# Build tool for Factorio Server Manager
#

build:
	# Build Linux release
	go build -o factorio-server-linux/factorio-server-manager src/*
	cp -r app/ factorio-server-linux/
	cp conf.json.example factorio-server-linux/conf.json
	zip -r $(HOME)/factorio-server-linux-x64.zip factorio-server-linux
	rm -rf factorio-server-linux
	# Build Windows release
	GOOS=windows GOARCH=386 go build -o factorio-server-windows/factorio-server-manager.exe src/*
	cp -r app/ factorio-server-windows/
	cp conf.json.example factorio-server-windows/conf.json
	zip -r $(HOME)/factorio-server-windows.zip factorio-server-windows
	rm -rf factorio-server-windows

