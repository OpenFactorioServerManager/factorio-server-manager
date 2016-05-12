# Build tool for Factorio Server Manager
#

build:
	go build -o $(HOME)/factorio-server/factorio-server-manager src/*
	cp -r app/ $(HOME)/factorio-server/
	cp conf.json.example $(HOME)/factorio-server/conf.json
