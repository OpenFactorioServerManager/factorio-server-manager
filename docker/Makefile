build:
	docker build -t factorio-server-manager .

run:
	docker run -d --name factorio-server-manager -v [yourpath]:/security -v [factorio-saves]:/opt/factorio/saves -v [factorio-mods]:/opt/factorio/mods -p 80:80 -p 443:443 -p 34197:34197/udp -e "ADMIN_PASSWORD=jqkSnQS4rA" factorio-server-manager
