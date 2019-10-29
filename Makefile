# Build tool for Factorio Server Manager

NODE_ENV:=production

#TODO add support for a mac build maybe?
UNAME := $(shell uname)
ifeq ($(UNAME), Linux)
	release := build/factorio-server-manager-linux.zip
else
	release := build/factorio-server-manager-windows.zip
endif

build: $(release)

$(shell mkdir -p build)
build/factorio-server-manager-%.zip: app/bundle factorio-server-manager-%
	@echo "Packaging Build - $@"
	@cp -r app/ factorio-server-manager/
	@cp conf.json.example factorio-server-manager/conf.json
	@zip -r $@ factorio-server-manager > /dev/null
	@rm -r factorio-server-manager

app/bundle:
	@echo "Building Frontend"
	@npm install && npm run build

factorio-server-manager-linux:
	@echo "Building Backend - Linux"
	@mkdir -p factorio-server-manager
	@cd src; \
	GOOS=linux GOARCH=amd64 go build -o ../factorio-server-manager/factorio-server-manager .

factorio-server-manager-windows:
	@echo "Building Backend - Windows"
	@mkdir -p factorio-server-manager
	@cd src; \
	GOOS=windows GOARCH=386 go build -o ../factorio-server-manager/factorio-server-manager.exe .

gen_release: build/factorio-server-manager-linux.zip build/factorio-server-manager-windows.zip
	@echo "Done"

clean:
	@echo "Cleaning"
	@rm -r build/
	@rm app/bundle.js
	@rm app/bundle.css
	@rm app/fonts/vendor
	@rm app/images/vendor
