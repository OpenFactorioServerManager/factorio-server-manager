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

build/factorio-server-manager-%.zip: clean app/bundle factorio-server-manager-%
	@mkdir -p build/
	@echo "Packaging Build - $@"
	@cp -r app/ factorio-server-manager/
	@cp conf.json.example factorio-server-manager/conf.json
	@zip -r $@ factorio-server-manager > /dev/null
	@rm -r factorio-server-manager/

app/bundle:
	@echo "Building Frontend"
	@npm install && npm run build

factorio-server-manager-linux:
	@echo "Building Backend - Linux"
	@mkdir -p factorio-server-manager
	@cd src; \
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ../factorio-server-manager/factorio-server-manager .

factorio-server-manager-windows:
	@echo "Building Backend - Windows"
	@mkdir -p factorio-server-manager
	@cd src; \
	GO111MODULE=on GOOS=windows GOARCH=amd64 CGO_ENABLED=0 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -ldflags="-extldflags=-static" -o ../factorio-server-manager/factorio-server-manager.exe .

gen_release: build/factorio-server-manager-linux.zip build/factorio-server-manager-windows.zip
	@echo "Done"

clean:
	@echo "Cleaning"
	@-rm -r build/
	@-rm app/bundle.js
	@-rm app/bundle.js.map
	@-rm app/style.css
	@-rm app/style.css.map
	@-rm -r app/fonts/vendor/
	@-rm -r app/images/vendor/
	@-rm -rf node_modules/
	@-rm -r pkg/
	@-rm -r factorio-server-manager
