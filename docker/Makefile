# Variables can be overridden by setting environment variables

FACTORIO_PATH ?= ~/.factorio
SECURITY_PATH ?= $(FACTORIO_PATH)/security
SAVES_PATH ?= $(FACTORIO_PATH)/saves
MODS_PATH ?= $(FACTORIO_PATH)/mods
PORT_FORWARD ?= -p 80:80 -p 443:443 -p 34197:34197/udp
FACTORIO_BRANCH ?= develop

build:
	docker build --build-arg FAC_BRANCH=$FACTORIO_BRANCH -f Dockerfile-build -t fsm-build .
	docker build -t factorio-server-manager .

logs:
	docker logs factorio-server -f

run:
	docker run -d --name factorio-server -v $(SECURITY_PATH):/security -v $(SAVES_PATH):/opt/factorio/saves -v $(MODS_PATH):/opt/factorio/mods $(PORT_FORWARD) factorio-server-manager

stop:
	docker stop factorio-server
	docker rm factorio-server

clean:
	docker rmi factorio-server-manager
	docker stop fsm-build
	docker rmi fsm-build
