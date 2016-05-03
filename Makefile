# Build tool for Factorio Server Manager
#

build:
	go build -o $(HOME)/factorio-server/factorio-server-manager
	cp -r app/ $(HOME)/factorio-server/
